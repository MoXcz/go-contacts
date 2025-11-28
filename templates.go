package main

import (
	"html/template"
	"path/filepath"
)

type pageData struct {
	Contacts   []Contact
	Contact    Contact
	SearchTerm string
	Page       int
	Errors     map[string]string
}

func add(x, y int) int { return x + y }
func sub(x, y int) int { return x - y }

var functions = template.FuncMap{
	"add": add,
	"sub": sub,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
