package session

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
