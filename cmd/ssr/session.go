package main

import "errors"

type SessionManager interface {
	Get(id string) *Session
	Create(id string) (*Session, error)
}

type SessionKeeper struct {
	sessions map[string]*Session
}

func NewSessionKeeper() SessionManager {
	return &SessionKeeper{
		sessions: make(map[string]*Session),
	}
}

func (sk *SessionKeeper) Get(id string) *Session {
	return sk.sessions[id]
}

func (sk *SessionKeeper) Create(id string) (*Session, error) {
	if _, exists := sk.sessions[id]; exists {
		return nil, errors.New("session already exists")
	}

	sk.sessions[id] = &Session{
		vals: make(map[string]string),
	}

	return sk.sessions[id], nil
}

type Session struct {
	vals map[string]string
}

func (s *Session) Flash(msg string) {
	s.vals["flash"] = msg
}

func (s *Session) GetFlash() string {
	msg := s.vals["flash"]
	delete(s.vals, "flash")
	return msg
}
