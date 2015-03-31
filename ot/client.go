package ot

type Client struct {
	Name      string    `json:"name"`
	Selection Selection `json:"selection"`
}
