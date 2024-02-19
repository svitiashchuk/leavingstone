package middleware

import (
	"leavingstone/internal/auth"
	"leavingstone/internal/handler"
	"leavingstone/internal/header"
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if auth.IsAuthenticated(r) {
			next(w, r)
			return
		}

		if header.IsHTMX(r) {
			handler.HTMXRedirect(w, r, "/login")
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}
