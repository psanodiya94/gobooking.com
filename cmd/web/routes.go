package main

import (
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/handlers"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)
	mux.Use(NoSurf)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/generals-quarters", handlers.Repo.Generals)

	mux.Get("/search-availability", handlers.Repo.GetAvailability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.JsonAvailability)

	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/make-reservation", handlers.Repo.GetReservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.GetShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)

	mux.Get("/user/logout", handlers.Repo.GetLogout)

	FileServer := http.FileServer(http.Dir(filepath.Join(".", "static")))
	mux.Handle("/static/*", http.StripPrefix("/static", FileServer))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.GetAdminDashboard)
		mux.Get("/reservations-all", handlers.Repo.GetAdminAllReservations)
		mux.Get("/reservations-new", handlers.Repo.GetAdminNewReservations)

		mux.Get("/process-reservations/{src}/{id}/do", handlers.Repo.GetAdminProcessReservation)
		mux.Get("/delete-reservations/{src}/{id}/do", handlers.Repo.GetAdminDeleteReservation)

		mux.Get("/reservations/{src}/{id}/show", handlers.Repo.GetAdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.PostAdminShowReservation)

		mux.Get("/reservations-calendar", handlers.Repo.GetAdminReservationsCalendar)
		mux.Post("/reservations-calendar", handlers.Repo.PostAdminReservationsCalendar)

	})

	return mux
}
