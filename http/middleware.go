package http

import "net/http"

func authMiddleware(endpointHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check session
		// if session exists, continue
		// if session does not exist, redirect to login
	}
}
