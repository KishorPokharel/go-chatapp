package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/KishorPokharel/chatapp/pkg/forms"
	"github.com/KishorPokharel/chatapp/pkg/models"
)

func (app *application) signupFormHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.html", &templateData{Form: forms.New(nil)})
}

func (app *application) loginFormHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.html", nil)
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	if form.Get("password") != form.Get("confirmpassword") {
		form.Errors.Add("confirmpassword", "Two passwords are not correct.")
	}
	form.Required("username", "password")
	form.MinLength("password", 8)
	if !form.Valid() {
		app.render(w, r, "register.html", &templateData{Form: form})
		return
	}

	user := &models.User{
		Username: strings.TrimSpace(form.Get("username")),
	}
	user.Password.Set(strings.TrimSpace(form.Get("password")))

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateUsername):
			form.Errors.Add("username", "Username already exists")
			app.render(w, r, "register.html", &templateData{Form: form})
			return
		default:
			app.serverError(w, r, err)
		}
		return
	}
	app.sessions.Put(r.Context(), "flash-success", "Successful. Please login here.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := app.models.Users.GetByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.sessions.Put(r.Context(), "flash-error", "Invalid Credentials")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		default:
			http.Error(w, "something went wrong", 500)
			return
		}
	}

	match, err := user.Password.Matches(password)
	if err != nil {
		app.logger.Println(err)
		return
	}
	if !match {
		app.sessions.Put(r.Context(), "flash-error", "Invalid Credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.sessions.Put(r.Context(), "chatapp:username", user.Username)
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	app.sessions.Remove(r.Context(), "chatapp:username")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessions.Exists(r.Context(), "chatapp:username")
}
