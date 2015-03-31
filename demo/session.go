package main

import "github.com/petejkim/ot.go/ot"

type Session struct {
	EventChan chan ConnEvent
	*ot.Session
}

func NewSession(document string) *Session {
	return &Session{
		make(chan ConnEvent),
		ot.NewSession(document),
	}
}

func (s *Session) HandleEvents() {
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
			username := data["username"]
			if username == nil {
				break
			}

			err := c.Send(&Event{"registered", c.ClientID})
			if err != nil {
				break
			}
			c.Broadcast(&Event{"join", map[string]interface{}{
				"client_id": c.ClientID,
				"username":  username,
			}})
		case "operation":
			data, ok := e.Data.([]interface{})
			if !ok {
				break
			}
			revf, ok := data[0].(float64)
			rev := int(revf)
			if !ok {
				break
			}
			ops, ok := data[1].([]interface{})
			if !ok {
				break
			}
			top, err := ot.Unmarshal(ops)
			if err != nil {
				break
			}
			top2, err := s.AddOperation(rev, top)
			if err != nil {
				break
			}

			err = c.Send(&Event{"ack", nil})
			if err != nil {
				break
			}
			c.Broadcast(&Event{"operation", []interface{}{c.ClientID, top2.Marshal()}})
		}
	}
}
