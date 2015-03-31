package session

import "github.com/nitrous-io/ot.go/ot/selection"

type Client struct {
	Name      string              `json:"name"`
	Selection selection.Selection `json:"selection"`
}
