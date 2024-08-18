package main

import "github.com/ga676005/snippetbox/internal/models"

type TemplateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
