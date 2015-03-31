package main

import (
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
)

var (
	nextClientID = 0
	connections  = map[*Connection]struct{}{}
)

type Event struct {
	Name string      `json:"e"`
	Data interface{} `json:"d,omitempty"`
}

type Connection struct {
	ClientID string
	Session  *Session
	Ws       *websocket.Conn
}

type ConnEvent struct {
	Conn *Connection
	*Event
}

func NewConnection(session *Session, ws *websocket.Conn) *Connection {
	return &Connection{Session: session, Ws: ws}
}

func (c *Connection) Handle() error {
	err := c.Send(&Event{"doc", map[string]interface{}{
		"document": c.Session.Document,
		"revision": len(c.Session.Operations),
		"clients":  c.Session.Clients,
	}})
	if err != nil {
		return err
	}

	RegisterConnection(c)
	c.ClientID = strconv.Itoa(nextClientID)
	c.Session.AddClient(c.ClientID)
	nextClientID++

	for {
		e, err := c.ReadEvent()
		if err != nil {
			break
		}

		c.Session.EventChan <- ConnEvent{c, e}
	}

	UnregisterConnection(c)
	if c.ClientID != "" {
		c.Session.RemoveClient(c.ClientID)
	}
	c.Broadcast(&Event{"quit", c.ClientID})

	return nil
}

func RegisterConnection(c *Connection) {
	connections[c] = struct{}{}
}

func UnregisterConnection(c *Connection) {
	delete(connections, c)
}

func (c *Connection) ReadEvent() (*Event, error) {
	_, msg, err := c.Ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	m := &Event{}
	if err = json.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Connection) Send(msg *Event) error {
	j, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err = c.Ws.WriteMessage(websocket.TextMessage, j); err != nil {
		return err
	}
	return nil
}

func (c *Connection) Broadcast(msg *Event) {
	for conn := range connections {
		if conn != c {
			conn.Send(msg)
		}
	}
}
