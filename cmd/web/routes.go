package main

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/psanodiya94/gobooking.com/pkg/config"
	"github.com/psanodiya94/gobooking.com/pkg/handlers"
)

/*---------------------------------------------------------------------------*/
func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)
	mux.Use(NoSurf)

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	FileServer := http.FileServer(http.Dir(filepath.Join(".", "static")))
	mux.Handle("/static/*", http.StripPrefix("/static", FileServer))

	return mux
}
