package auth

import (
	"context"
	"leavingstone/internal/model"
	"leavingstone/internal/sqlite"
	"net/http"
)

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

const AuthCookie = "auth_token"

const UserIDContextKey = contextKey("userID")

const IsAuthenticatedContextKey = contextKey("isAuthenticated")

type Authenticator struct {
	us *sqlite.UserService
}

func New(us *sqlite.UserService) *Authenticator {
	return &Authenticator{
		us: us,
	}
}

func IsAuthenticated(r *http.Request) bool {
	v := r.Context().Value(isAuthenticatedContextKey)
	if v == nil {
		return false
	}

	isAuthenticated := v.(bool)
	return isAuthenticated
}

func (auth *Authenticator) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var u *model.User
		for _, c := range r.Cookies() {
			if c.Name == AuthCookie && c.Value != "" {
				u, err := auth.us.FindByToken(c.Value)
				if err != nil {
					panic(err)
				}

				if u != nil {
					break
				}
			}
		}

		if u != nil {
			ctx := r.Context()
			ctx = context.WithValue(ctx, IsAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, UserIDContextKey, u.ID)

			next(w, r.WithContext(ctx))
			return
		}

		next(w, r)
	}
}
