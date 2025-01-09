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

var functions = template.FuncMap{}

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var templatePath = "./../../templates"

func TestMain(m *testing.M) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

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

	FileServer := http.FileServer(http.Dir(filepath.Join(".", "static")))
	mux.Handle("/static/*", http.StripPrefix("/static", FileServer))

	return mux
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

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
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", templatePath))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", templatePath))
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
