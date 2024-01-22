package main

import (
	"context"
	"net/http"
)

func (app *App) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			next(w, r)
			return
		}

		if isHTMX(r) {
			app.htmxRedirect(w, r, "/login")
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

func (app *App) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := app.auth.authenticate(r)
		if u != nil {
			ctx := r.Context()
			ctx = context.WithValue(ctx, isAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, userIDContextKey, u.ID)

			next(w, r.WithContext(ctx))
			return
		}

		next(w, r)
	}
}
