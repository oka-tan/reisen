package main

import (
	"io"
	"log"

	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
)

type Templater struct {
	layout    *mustache.Template
	templates map[string]*mustache.Template
}

func mustCompile(filename string, partials mustache.PartialProvider) *mustache.Template {
	template, err := mustache.ParseFilePartials(filename, partials)

	if err != nil {
		log.Fatalln(err)
	}

	return template
}

func NewTemplater() *Templater {
	fileProvider := &mustache.FileProvider{
		Paths:      []string{"templates/"},
		Extensions: []string{".html.mustache-partial"},
	}

	layout := mustCompile("templates/layout.html.mustache", fileProvider)

	templates := map[string]*mustache.Template{
		"index":                         mustCompile("templates/index.html.mustache", fileProvider),
		"board-catalog-variant":         mustCompile("templates/board-catalog-variant.html.mustache", fileProvider),
		"board-error":                   mustCompile("templates/board-error.html.mustache", fileProvider),
		"board-empty":                   mustCompile("templates/board-empty.html.mustache", fileProvider),
		"board-search":                  mustCompile("templates/board-search.html.mustache", fileProvider),
		"board-search-thread-not-found": mustCompile("templates/board-search-thread-not-found.html.mustache", fileProvider),
		"board-search-bad-request":      mustCompile("templates/board-search-bad-request.html.mustache", fileProvider),
		"board-search-no-results":       mustCompile("templates/board-search-no-results.html.mustache", fileProvider),
		"board-search-server-error":     mustCompile("templates/board-search-server-error.html.mustache", fileProvider),
		"board-thread":                  mustCompile("templates/board-thread.html.mustache", fileProvider),
		"board-thread-error":            mustCompile("templates/board-thread-error.html.mustache", fileProvider),
		"board-thread-not-found":        mustCompile("templates/board-thread-not-found.html.mustache", fileProvider),
		"board-post-not-found":          mustCompile("templates/board-post-not-found.html.mustache", fileProvider),
		"board-view-same":               mustCompile("templates/board-view-same.html.mustache", fileProvider),
		"board-view-same-error":         mustCompile("templates/board-view-same-error.html.mustache", fileProvider),
		"board-view-same-empty":         mustCompile("templates/board-view-same-empty.html.mustache", fileProvider),
		"contact":                       mustCompile("templates/contact.html.mustache", fileProvider),
		"search-reference":              mustCompile("templates/search-reference.html.mustache", fileProvider),
	}

	return &Templater{
		layout:    layout,
		templates: templates,
	}
}

//Render implements a method echo needs for template rendering.
func (t *Templater) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates[name].FRenderInLayout(w, t.layout, data)
}
