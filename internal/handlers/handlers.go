package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"github.com/psanodiya94/gobooking.com/internal/render"
	"log"
	"net/http"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the homepage handler
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderingTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some business logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// get the user's IP address from the session
	stringMap["remote_ip"] = repo.App.Session.GetString(r.Context(), "remote_ip")

	render.RenderingTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals is the generals page handler
func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderingTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the majors page handler
func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderingTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Reservation is the reservation page handler
func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderingTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Contact is the contact page handler
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderingTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// GetAvailability checks the availability of rooms
func (repo *Repository) GetAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderingTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability checks the availability of rooms
func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	_, _ = w.Write([]byte(fmt.Sprintf("Start Date is %s && End Date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// JsonAvailability checks the availability of rooms
func (repo *Repository) JsonAvailability(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}
