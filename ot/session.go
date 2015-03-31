package ot

type Session struct {
	Document string
}

func NewSession(document string) *Session {
	return &Session{Document: document}
}
