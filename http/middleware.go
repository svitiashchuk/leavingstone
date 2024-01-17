package http

import "net/http"

func (app *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !app.auth.isAuthenticated(r) {
			app.htmxRedirect(w, r, "/login")
			w.Write([]byte("You are not authorized to access this page"))
			return
		}

		next(w, r)
	}
}
