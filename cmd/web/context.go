package main

import (
	"context"
	"net/http"

	"github.com/KishorPokharel/chatapp/pkg/models"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		panic("missing user context")
	}
	return user
}
