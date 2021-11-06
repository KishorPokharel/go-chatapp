package main

import (
	"bytes"
	"net/http"
)

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.Flash = map[string]string{
		"success": app.sessions.PopString(r.Context(), "flash-success"),
		"error":   app.sessions.PopString(r.Context(), "flash-error"),
	}
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

func (app *application) render(
	w http.ResponseWriter,
	r *http.Request,
	filename string,
	data *templateData,
) {
	tmpl, ok := app.templateCache[filename]
	if !ok {
		app.logger.Panicf("template %s not in cache", filename)
	}
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, app.addDefaultData(data, r))
	// TODO: needs proper error handling
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

/*
func (app *application) notFoundError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.render(w, r, "404.html", nil)
}
*/
func (app *application) serverError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.logger.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	app.render(w, r, "500.html", nil)
}
