package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/KishorPokharel/chatapp/pkg/forms"
	"github.com/KishorPokharel/chatapp/pkg/models"
)

type templateData struct {
	CSRFToken       string
	Flash           map[string]string
	Form            *forms.Form
	IsAuthenticated bool
	User            *models.User
	Messages        []*models.Message
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("01/02/2006 15:04 PM")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New("base").Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
