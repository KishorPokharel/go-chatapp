package main

import "net/http"

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		username := app.sessions.GetString(r.Context(), "chatapp:username")

		user, err := app.models.Users.GetByUsername(username)
		if err != nil {
			http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
			return
		}
		r = app.contextSetUser(r, user)
		next(w, r)
	}
}

func (app *application) requireAnonymous(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
