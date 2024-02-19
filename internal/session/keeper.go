package session

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
)

type contextKey string

const SessionContextKey = contextKey("session")

type Keeper struct {
	sessions map[string]*Session
}

func NewKeeper() *Keeper {
	return &Keeper{
		sessions: make(map[string]*Session),
	}
}

func (sk *Keeper) Get(id string) *Session {
	return sk.sessions[id]
}

func (sk *Keeper) Create(id string) (*Session, error) {
	if _, exists := sk.sessions[id]; exists {
		return nil, errors.New("session already exists")
	}

	sk.sessions[id] = &Session{
		vals: make(map[string]string),
	}

	return sk.sessions[id], nil
}

func (sm *Keeper) Provide(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Cookies()

		for _, cookie := range c {
			if cookie.Name == "session_id" {
				s := sm.Get(cookie.Value)
				if s != nil {
					ctx := r.Context()
					ctx = context.WithValue(ctx, SessionContextKey, cookie.Value)

					next(w, r.WithContext(ctx))
					return
				}
			}
		}

		sessionID := strconv.Itoa(rand.Intn(1000000000))
		_, err := sm.Create(sessionID)
		if err != nil {
			panic(err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
			Path:  "/",
		})

		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionContextKey, sessionID)

		next(w, r.WithContext(ctx))
	}
}
