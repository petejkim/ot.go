package main

import (
	"strconv"
	"sync"

	"github.com/nitrous-io/ot.go/ot/operation"
	"github.com/nitrous-io/ot.go/ot/selection"
	"github.com/nitrous-io/ot.go/ot/session"
)

type Session struct {
	nextConnID  int
	Connections map[*Connection]struct{}

	EventChan chan ConnEvent

	lock sync.Mutex

	*session.Session
}

func NewSession(document string) *Session {
	return &Session{
		Connections: map[*Connection]struct{}{},
		EventChan:   make(chan ConnEvent),
		Session:     session.New(document),
	}
}

func (s *Session) RegisterConnection(c *Connection) {
	s.lock.Lock()
	id := strconv.Itoa(s.nextConnID)
	c.ID = id
	s.nextConnID++
	s.Connections[c] = struct{}{}
	s.AddClient(c.ID)
	s.lock.Unlock()
}

func (s *Session) UnregisterConnection(c *Connection) {
	s.lock.Lock()
	delete(s.Connections, c)
	if c.ID != "" {
		s.RemoveClient(c.ID)
	}
	s.lock.Unlock()
}

func (s *Session) HandleEvents() {
	// this method should run in a single go routine
	for {
		e, ok := <-s.EventChan
		if !ok {
			return
		}

		c := e.Conn
		switch e.Name {
		case "join":
			data, ok := e.Data.(map[string]interface{})
			if !ok {
				break
			}
			username, ok := data["username"].(string)
			if !ok || username == "" {
				break
			}

			s.SetName(c.ID, username)

			err := c.Send(&Event{"registered", c.ID})
			if err != nil {
				break
			}
			c.Broadcast(&Event{"join", map[string]interface{}{
				"client_id": c.ID,
				"username":  username,
			}})
		case "op":
			// data: [revision, ops, selection?]
			data, ok := e.Data.([]interface{})
			if !ok {
				break
			}
			if len(data) < 2 {
				break
			}
			// revision
			revf, ok := data[0].(float64)
			rev := int(revf)
			if !ok {
				break
			}
			// ops
			ops, ok := data[1].([]interface{})
			if !ok {
				break
			}
			top, err := operation.Unmarshal(ops)
			if err != nil {
				break
			}
			// selection (optional)
			if len(data) >= 3 {
				selm, ok := data[2].(map[string]interface{})
				if !ok {
					break
				}
				sel, err := selection.Unmarshal(selm)
				if err != nil {
					break
				}
				top.Meta = sel
			}

			top2, err := s.AddOperation(rev, top)
			if err != nil {
				break
			}

			err = c.Send(&Event{"ok", nil})
			if err != nil {
				break
			}

			if sel, ok := top2.Meta.(*selection.Selection); ok {
				c.Broadcast(&Event{"op", []interface{}{c.ID, top2.Marshal(), sel.Marshal()}})
			} else {
				c.Broadcast(&Event{"op", []interface{}{c.ID, top2.Marshal()}})
			}
		case "sel":
			data, ok := e.Data.(map[string]interface{})
			if !ok {
				break
			}
			sel, err := selection.Unmarshal(data)
			if err != nil {
				break
			}
			c.Broadcast(&Event{"sel", []interface{}{c.ID, sel.Marshal()}})
		}
	}
}
