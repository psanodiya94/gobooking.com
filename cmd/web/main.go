package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/driver"
	"github.com/psanodiya94/gobooking.com/internal/handlers"
	"github.com/psanodiya94/gobooking.com/internal/helpers"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"github.com/psanodiya94/gobooking.com/internal/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	log.Println("Starting mail listener")

	listenForMail()

	log.Println("Starting application on port", port)

	server := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	log.Fatal(server.ListenAndServe())
}

func run() (*driver.DataBase, error) {
	// what will be put in the session
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "application is in production")
	useCache := flag.Bool("cache", true, "use template cache")
	dbHost := flag.String("dbhost", "localhost", "database name")
	dbName := flag.String("dbname", "", "database name")
	dbUser := flag.String("dbuser", "", "database user")
	dbPass := flag.String("dbpass", "", "database password")
	dbPort := flag.String("dbport", "5432", "database port")
	dbSSL := flag.String("dbssl", "disable", "database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing Required Flags")
		os.Exit(1)
	}

	// setup mail channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// initialize loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	conn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL,
	)
	db, err := driver.ConnectSql(conn)
	if err != nil {
		log.Fatal("Can't connect to database! Dying...")
		return nil, err
	}
	log.Println("Connected to database!")

	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tmplCache

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	render.NewRenderer(&app)

	return db, nil
}
