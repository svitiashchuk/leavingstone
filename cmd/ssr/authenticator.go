package main

import (
	"leavingstone/sqlite"
	"net/http"
)

const AuthCookie = "auth_token"

type Authenticator struct {
	us *sqlite.UserService
}

func (auth *Authenticator) isAuthenticated(r *http.Request) bool {
	for _, c := range r.Cookies() {
		if c.Name == AuthCookie {
			token := c.Value

			if token == "" {
				return false
			}

			u, err := auth.us.FindByToken(token)
			if err != nil {
				panic(err)
			}

			return u != nil
		}
	}

	return false
}
