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

func mustCompile(filename string) *mustache.Template {
	template, err := mustache.ParseFile(filename)

	if err != nil {
		log.Fatalln(err)
	}

	return template
}

func NewTemplater() *Templater {
	layout := mustCompile("templates/layout.html.mustache")

	templates := map[string]*mustache.Template{
		"index":                         mustCompile("templates/index.html.mustache"),
		"board":                         mustCompile("templates/board.html.mustache"),
		"board-error":                   mustCompile("templates/board-error.html.mustache"),
		"board-empty":                   mustCompile("templates/board-empty.html.mustache"),
		"board-search":                  mustCompile("templates/board-search.html.mustache"),
		"board-search-thread-not-found": mustCompile("templates/board-search-thread-not-found.html.mustache"),
		"board-search-bad-request":      mustCompile("templates/board-search-bad-request.html.mustache"),
		"board-search-no-results":       mustCompile("templates/board-search-no-results.html.mustache"),
		"board-search-server-error":     mustCompile("templates/board-search-server-error.html.mustache"),
		"board-thread":                  mustCompile("templates/board-thread.html.mustache"),
		"board-thread-error":            mustCompile("templates/board-thread-error.html.mustache"),
		"board-thread-thread-not-found": mustCompile("templates/board-thread-not-found.html.mustache"),
		"contact":                       mustCompile("templates/contact.html.mustache"),
		"search-reference":              mustCompile("templates/search-reference.html.mustache"),
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
