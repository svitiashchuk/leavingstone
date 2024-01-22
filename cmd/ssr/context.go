package main

import "net/http"

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

func (app *App) isAuthenticated(r *http.Request) bool {
	isAuthenticated := r.Context().Value(isAuthenticatedContextKey).(bool)
	return isAuthenticated
}
