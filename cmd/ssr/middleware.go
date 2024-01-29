package main

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
)

func (app *App) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
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

func (app *App) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := app.auth.authenticate(r)
		if u != nil {
			ctx := r.Context()
			ctx = context.WithValue(ctx, isAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, userIDContextKey, u.ID)

			next(w, r.WithContext(ctx))
			return
		}

		next(w, r)
	}
}

func (app *App) session(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Cookies()

		for _, cookie := range c {
			if cookie.Name == "session_id" {
				s := app.sm.Get(cookie.Value)
				if s != nil {
					ctx := r.Context()
					ctx = context.WithValue(ctx, sessionContextKey, cookie.Value)

					next(w, r.WithContext(ctx))
					return
				}
			}
		}

		sessionID := strconv.Itoa(rand.Intn(1000000000))
		_, err := app.sm.Create(sessionID)
		if err != nil {
			panic(err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
			Path:  "/",
		})

		ctx := r.Context()
		ctx = context.WithValue(ctx, sessionContextKey, sessionID)

		next(w, r.WithContext(ctx))
	}
}
