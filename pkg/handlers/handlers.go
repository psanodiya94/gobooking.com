package handlers

import (
	"net/http"

	"github.com/psanodiya94/gobooking.com/pkg/config"
	"github.com/psanodiya94/gobooking.com/pkg/models"
	"github.com/psanodiya94/gobooking.com/pkg/render"
)

/*---------------------------------------------------------------------------*/
// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

/*---------------------------------------------------------------------------*/
// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

/*---------------------------------------------------------------------------*/
// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

/*---------------------------------------------------------------------------*/
// Home is the homepage handler
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

/*---------------------------------------------------------------------------*/
// About is the about page handler
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some business logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// get the user's IP address from the session
	stringMap["remote_ip"] = repo.App.Session.GetString(r.Context(), "remote_ip")

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
