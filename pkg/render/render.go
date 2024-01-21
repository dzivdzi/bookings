package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dzivdzi/bookings/pkg/config"
	"github.com/dzivdzi/bookings/pkg/models"
)

var tc = make(map[string]*template.Template)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// When we dicide what data should be on EVERY PAGE we can add it here
// Currently, it just takes the data from models.TemplateData and is returning it
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// Get the template cache from the app config
	// takes this from main and uses it because we don't wanna create the template cache always
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get requested template from cache

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// Render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// Declare a map pointing towards the template - this is what we can actually render
	// Produces a safe HTML document fragment
	myCache := map[string]*template.Template{}

	// Get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Range through all files ending with *.page.tmpl
	// first underscore is the index, we don't care about that
	// page is a var declared to store the single return of the
	// pages map above - lamens terms -> for every page in pages
	for _, page := range pages {
		// This func (Base) slices the string assuming that the separator is "/" and saves ONLY the last entry after the last "/"
		// So if the full path is /template/index.page.tmpl, it will only take the index.page.tmpl as stated in the func
		name := filepath.Base(page)
		// ts is just a pointer to template.Template - what we do with it is we give it a name with template.NEW() and then we parse the file page
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// This serves the same function as above but now we are looking for files that end with .layout.tmpl
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		// finally, here we add the final result to the dict(map)
		myCache[name] = ts
	}

	return myCache, nil

}
