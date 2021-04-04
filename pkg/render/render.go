package render

import (
	"bytes"
	"fmt"
	"github.com/usmanzaheer1995/bed-and-breakfast/pkg/config"
	"github.com/usmanzaheer1995/bed-and-breakfast/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the templates package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// this gets a list of all files ending with page.tmpl, and stores
	// it in a slice of strings called pages
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// now we loop through the slice of strings, which has two
	// entries: "home.page.tmpl" and "about.page.tmpl"
	for _, page := range pages {
		// the first time through, name = "home.page.tmpl"
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// here, we are checking to see if there are any files at all that
		// end with layout.tmpl. THere is only one, but if there were more
		// than one, we we get them all and store them in a slice of strings
		// named matches
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// if the length of matches is > 0, we go through the slice
		// and parse all of the layouts available to us. We might not use
		// any of them in this iteration through the loop, but if the current
		// template we are working on (home.page.tmpl the first time through)
		// does use a layout, we need to have it available to us before we add it
		// to our template set
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		// the first time through, name is still home.page.tmpl
		// we never add anything with *.layout.tmpl to the template set;
		// we just use the layout to create a page which depends on it.
		// now, we add the template, complete any associated layouts, to our
		// template set
		myCache[name] = ts
	}
	return myCache, nil
}
