package main

import "net/http"

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

const userIDContextKey = contextKey("userID")

func (app *App) isAuthenticated(r *http.Request) bool {
	isAuthenticated := r.Context().Value(isAuthenticatedContextKey).(bool)
	return isAuthenticated
}

func (app *App) userID(r *http.Request) int {
	return r.Context().Value(userIDContextKey).(int)
}
