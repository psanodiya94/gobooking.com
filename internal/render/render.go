package render

import (
	"bytes"
	"github.com/justinas/nosurf"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderingTemplate renders a template using html template
func RenderingTemplate(w http.ResponseWriter, r *http.Request, tmpl string, tmplData *models.TemplateData) {
	var tmplCache map[string]*template.Template

	if app.UseCache {
		// get the t cache from the app config
		tmplCache = app.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	// get requested t from cache
	t, ok := tmplCache[tmpl]
	if !ok {
		log.Fatal("Could not get t from t cache")
	}

	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, r)

	_ = t.Execute(buf, tmplData)

	// render the t
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing t to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all the files name in the templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return cache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
