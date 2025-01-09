package render

import (
	"github.com/psanodiya94/gobooking.com/internal/models"
	"net/http"
	"testing"
)

// TestAddDefaultData tests AddDefaultData function
func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	resp, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	session.Put(resp.Context(), "flash", "flash message")
	result := AddDefaultData(&td, resp)
	if result.Flash != "flash message" {
		t.Error("flash value of 'flash message' not found in session")
	}
}

// TestTemplate tests for Template function
func TestTemplate(t *testing.T) {
	templatePath = "./../../templates"
	templCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = templCache

	session, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	var writer writerResponse

	err = Template(&writer, session, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&writer, session, "non.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered a template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	templatePath = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSessionData() (*http.Request, error) {
	resp, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := resp.Context()
	ctx, _ = session.Load(ctx, resp.Header.Get("X-Session"))

	resp = resp.WithContext(ctx)

	return resp, nil
}
