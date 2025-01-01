package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/psanodiya94/gobooking.com/pkg/config"
	"github.com/psanodiya94/gobooking.com/pkg/models"
)

/*---------------------------------------------------------------------------*/
/*---------------------------------Option 1----------------------------------*/
// RenderTemplate renders a template using html template
// func RenderTemplate(w http.ResponseWriter, tmpl string) {
// 	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
// 	err := parsedTemplate.Execute(w, nil)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

/*---------------------------------------------------------------------------*/
/*---------------------------------Option 2----------------------------------*/
// var cache = make(map[string]*template.Template)

/*---------------------------------------------------------------------------*/
// func RenderTemplate(w http.ResponseWriter, tmpl string) {
// 	var tmplCache *template.Template
// 	var err error
// 	_, inMap := cache[tmpl]
// 	if !inMap {
// 		// we need to add the template to the cache
// 		log.Println("creating template and adding to cache")
// 		err = createTemplateCache(tmpl)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		// we are using the template from the cache
// 		log.Println("Using template from cache")
// 	}

// 	tmplCache = cache[tmpl]

// 	err = tmplCache.Execute(w, nil)
// 	if err != nil {
// 		log.Printf("Error executing template: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

/*---------------------------------------------------------------------------*/
// func createTemplateCache(tmpl string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", tmpl),
// 		"./templates/base.layout.tmpl",
// 	}

// 	// this is a slice of strings
// 	parsedTemplate, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	// add template to cache
// 	cache[tmpl] = parsedTemplate

// 	return nil
// }

/*---------------------------------------------------------------------------*/
/*------------------------(Recomended)-Option 3------------------------------*/
var app *config.AppConfig

/*---------------------------------------------------------------------------*/
// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

/*---------------------------------------------------------------------------*/
// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

/*---------------------------------------------------------------------------*/
// RenderTemplate renders a template using html template
func RenderTemplate(w http.ResponseWriter, tmpl string, tmplData *models.TemplateData) {
	var tmplCache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tmplCache = app.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	template, ok := tmplCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData)

	_ = template.Execute(buf, tmplData)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}
}

/*---------------------------------------------------------------------------*/
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all of the files name in the templates folder
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

/*---------------------------------------------------------------------------*/
