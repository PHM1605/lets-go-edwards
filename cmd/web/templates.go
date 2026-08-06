// To hold all data that needs to be passed around (not only Snippet struct)
package main

import (
	"html/template"
	"lets-go-edwards/internal/models"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

// NOTE: to pass Functions to Template
// 1/ define a normal function
// 2/ create a "template.FuncMap"
// 3/ use "template.New("xxx").Funcs(<FuncMap>)"
// 4/ use that function in Template with {{humanDate .Created}} or PIPELINING {{.Created | humanDate}}
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// map "name of template" and that Template
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page) // e.g. home.html
		// parse base template
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// parse all partials in terms of a Glob pattern
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")

		// parse current page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add current page to cache
		cache[name] = ts
	}
	return cache, nil
}
