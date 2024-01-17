package http

import (
	"leavingstone/sqlite"
	"net/http"
)

const HTTPHeaderAuthToken = "X-Auth-Token"

type Authenticator struct {
	us *sqlite.UserService
}

func (auth *Authenticator) isAuthenticated(r *http.Request) bool {
	token := r.Header.Get(HTTPHeaderAuthToken)
	if token == "" {
		return false
	}

	u, err := auth.us.FindByToken(token)
	if err != nil {
		panic(err)
	}

	return u != nil
}
