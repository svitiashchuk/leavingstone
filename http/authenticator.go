package http

import "net/http"

const HTTPHeaderSessionID = "X-App-Session-ID"

func (app *Server) isAuthenticated(r *http.Request) bool {
	sID := r.Header.Get(HTTPHeaderSessionID)
	if sID == "" {
		return false
	}

	// Validate the auth header using the session manager
	return app.sm.Get(sID) != nil
}
