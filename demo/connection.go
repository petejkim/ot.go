package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var nextClientID = 0
var connections = map[*Connection]struct{}{}
var lock = &sync.Mutex{}

type Message struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

type Connection struct {
	ClientID string
	Session  *Session
	Ws       *websocket.Conn
}

func NewConnection(session *Session, ws *websocket.Conn) *Connection {
	return &Connection{Session: session, Ws: ws}
}

func (c *Connection) Handle() error {
	err := c.Send(&Message{"doc", map[string]interface{}{
		"document": session.Document,
		"revision": len(session.Operations),
		"clients":  session.Clients,
	}})
	if err != nil {
		return err
	}

	RegisterConnection(c)

	for {
		m, err := c.Read()
		if err != nil {
			break
		}

		switch m.Event {
		case "join":
			if c.ClientID != "" {
				break
			}
			username := m.Data["username"]
			if username == nil {
				break
			}
			c.ClientID = fmt.Sprintf("%d:%s", nextClientID, username)
			nextClientID++
			session.AddClient(c.ClientID)

			err = c.Send(&Message{"client_id", map[string]interface{}{
				"client_id": c.ClientID,
			}})
			if err != nil {
				break
			}
			Broadcast(&Message{"join", map[string]interface{}{
				"client_id": c.ClientID,
			}})
		}
	}

	UnregisterConnection(c)
	if c.ClientID != "" {
		session.RemoveClient(c.ClientID)
	}
	Broadcast(&Message{"quit", map[string]interface{}{
		"client_id": c.ClientID,
	}})

	return nil
}

func RegisterConnection(c *Connection) {
	lock.Lock()
	connections[c] = struct{}{}
	lock.Unlock()
}

func UnregisterConnection(c *Connection) {
	lock.Lock()
	delete(connections, c)
	lock.Unlock()
}

func (c *Connection) Read() (*Message, error) {
	_, msg, err := c.Ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	m := &Message{}
	if err = json.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Connection) Send(msg *Message) error {
	j, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err = c.Ws.WriteMessage(websocket.TextMessage, j); err != nil {
		return err
	}
	return nil
}

func Broadcast(msg *Message) {
	for c := range connections {
		c.Send(msg)
	}
}
