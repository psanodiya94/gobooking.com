package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

var functions = template.FuncMap{
	"readableDate": ReadableDate,
	"formatDate":   FormatDate,
	"iterate":      Iterate,
	"add":          Add,
}

var app *config.AppConfig
var templatePath = "./templates"

// Add adds a & b and returns
func Add(a, b int) int {
	return a + b
}

// Iterate returns a slice of int, starting at 1, going to count
func Iterate(count int) []int {
	// Create a slice of int
	var items []int
	for i := 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// ReadableDate returns time in YYYY-MM-DD format
func ReadableDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate formate date with string
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuth = true
	}

	return td
}

// Template renders a template using html template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, tmplData *models.TemplateData) error {
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
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, r)

	_ = t.Execute(buf, tmplData)

	// render the t
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing t to browser", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
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
