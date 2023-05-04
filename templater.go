package main

import (
	"io"
	"log"

	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
)

type templater struct {
	layout    *mustache.Template
	templates map[string]templateAndLayout
}

type templateAndLayout struct {
	layout   *mustache.Template
	template *mustache.Template
}

func mustCompile(filename string, partials mustache.PartialProvider) *mustache.Template {
	template, err := mustache.ParseFilePartials(filename, partials)

	if err != nil {
		log.Fatalln(err)
	}

	return template
}

func newTemplater() *templater {
	fileProvider := &mustache.FileProvider{
		Paths:      []string{"templates/"},
		Extensions: []string{".html.mustache-partial"},
	}

	layout := mustCompile("templates/layout.html.mustache", fileProvider)

	templates := map[string]templateAndLayout{
		"index": {
			template: mustCompile("templates/index.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-catalog-variant": {
			template: mustCompile("templates/board-catalog-variant.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-error": {
			template: mustCompile("templates/board-error.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-empty": {
			template: mustCompile("templates/board-empty.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-search": {
			template: mustCompile("templates/board-search.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-search-thread-not-found": {
			template: mustCompile("templates/board-search-thread-not-found.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-search-bad-request": {
			template: mustCompile("templates/board-search-bad-request.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-search-no-results": {
			template: mustCompile("templates/board-search-no-results.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-search-server-error": {
			template: mustCompile("templates/board-search-server-error.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-thread": {
			template: mustCompile("templates/board-thread.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-thread-error": {
			template: mustCompile("templates/board-thread-error.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-thread-not-found": {
			template: mustCompile("templates/board-thread-not-found.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-post-not-found": {
			template: mustCompile("templates/board-post-not-found.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-view-same": {
			template: mustCompile("templates/board-view-same.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-view-same-error": {
			template: mustCompile("templates/board-view-same-error.html.mustache", fileProvider),
			layout:   layout,
		},
		"board-view-same-empty": {
			template: mustCompile("templates/board-view-same-empty.html.mustache", fileProvider),
			layout:   layout,
		},
		"about": {
			template: mustCompile("templates/about.html.mustache", fileProvider),
			layout:   layout,
		},

		"report": {
			template: mustCompile("templates/report.html.mustache", fileProvider),
			layout:   nil,
		},
		"report-result": {
			template: mustCompile("templates/report-result.html.mustache", fileProvider),
			layout:   nil,
		},
	}

	return &templater{
		layout:    layout,
		templates: templates,
	}
}

// Render implements a method echo needs for template rendering.
func (t *templater) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templateAndLayout := t.templates[name]

	if templateAndLayout.layout == nil {
		return templateAndLayout.template.FRender(w, data)
	}

	return templateAndLayout.template.FRenderInLayout(w, templateAndLayout.layout, data)
}
