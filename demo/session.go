package main

import "github.com/petejkim/ot.go/ot"

type Session struct {
	Clients []string
	*ot.Session
}

func NewSession(document string) *Session {
	return &Session{
		[]string{},
		ot.NewSession(document),
	}
}

func (s *Session) AddClient(id string) {
	s.Clients = append(s.Clients, id)
}

func (s *Session) RemoveClient(id string) {
	i := -1
	for j, u := range s.Clients {
		if u == id {
			i = j
			break
		}
	}
	if i == -1 {
		return
	}
	s.Clients = append(s.Clients[:i], s.Clients[i+1:]...)
}
