package ot

type Operation struct {
}

type Session struct {
	Document   string
	Operations []*Operation
}

func NewSession(document string) *Session {
	return &Session{Document: document}
}
