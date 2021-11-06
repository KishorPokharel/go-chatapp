package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	staticHandler := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticHandler))

	r.Get("/register", app.requireAnonymous(app.signupFormHandler))
	r.Get("/login", app.requireAnonymous(app.loginFormHandler))
	r.Post("/register", app.requireAnonymous(app.signupHandler))
	r.Post("/login", app.requireAnonymous(app.loginHandler))
	r.Post("/logout", app.requireAuthentication(app.logoutHandler))

	r.Get("/chat", app.requireAuthentication(app.chatHandler))
	r.HandleFunc("/room", app.requireAuthentication(app.roomHandler))

	return app.sessions.LoadAndSave(r)
}
