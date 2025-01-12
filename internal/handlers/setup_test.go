package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"github.com/psanodiya94/gobooking.com/internal/render"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"text/template"
	"time"
)

var functions = template.FuncMap{
	"readableDate": render.ReadableDate,
	"formatDate":   render.FormatDate,
	"iterate":      render.Iterate,
	"add":          render.Add,
}

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var templatePath = "./../../templates"

func TestMain(m *testing.M) {
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// initialize loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	tmplCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tmplCache
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/generals-quarters", Repo.Generals)

	mux.Get("/search-availability", Repo.GetAvailability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.JsonAvailability)

	mux.Get("/make-reservation", Repo.GetReservation)
	mux.Post("/make-reservation", Repo.PostReservation)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/user/login", Repo.GetShowLogin)
	mux.Post("/user/login", Repo.PostShowLogin)

	mux.Get("/user/logout", Repo.GetLogout)

	mux.Get("/admin/dashboard", Repo.GetAdminDashboard)
	mux.Get("/admin/reservations-all", Repo.GetAdminAllReservations)
	mux.Get("/admin/reservations-new", Repo.GetAdminNewReservations)

	mux.Get("/admin/process-reservations/{src}/{id}/do", Repo.GetAdminProcessReservation)
	mux.Get("/admin/delete-reservations/{src}/{id}/do", Repo.GetAdminDeleteReservation)

	mux.Get("/admin/reservations/{src}/{id}/show", Repo.GetAdminShowReservation)
	mux.Post("/admin/reservations/{src}/{id}", Repo.PostAdminShowReservation)

	mux.Get("/admin/reservations-calendar", Repo.GetAdminReservationsCalendar)
	mux.Post("/admin/reservations-calendar", Repo.PostAdminReservationsCalendar)

	FileServer := http.FileServer(http.Dir(filepath.Join(".", "static")))
	mux.Handle("/static/*", http.StripPrefix("/static", FileServer))

	return mux
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all the files name in the templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", templatePath))
	if err != nil {
		return cache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Println(err)
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", templatePath))
		if err != nil {
			log.Println(err)
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", templatePath))
			if err != nil {
				log.Println(err)
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
