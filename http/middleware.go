package http

import "net/http"

func (app *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}
