package session

type Manager interface {
	Get(id string) *Session
	Create(id string) (*Session, error)
}
