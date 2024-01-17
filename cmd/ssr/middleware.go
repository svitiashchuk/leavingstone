package main

import "net/http"

func (app *App) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.auth.isAuthenticated(r) {
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
