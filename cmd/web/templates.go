package main

import (
	"path/filepath"
	"text/template"

	"github.com/ga676005/snippetbox/internal/models"
)

type TemplateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable
		name := filepath.Base(page)

		// base
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return cache, err
		}

		// all partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return cache, err
		}

		// page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return cache, err
		}

		cache[name] = ts
	}

	return cache, nil
}
