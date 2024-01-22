package main

import (
	"leavingstone"
	"leavingstone/sqlite"
	"net/http"
)

const AuthCookie = "auth_token"

type Authenticator struct {
	us *sqlite.UserService
}

func (auth *Authenticator) authenticate(r *http.Request) *leavingstone.User {
	for _, c := range r.Cookies() {
		if c.Name == AuthCookie {
			token := c.Value

			if token == "" {
				return nil
			}

			u, err := auth.us.FindByToken(token)
			if err != nil {
				panic(err)
			}

			return u
		}
	}

	return nil
}
