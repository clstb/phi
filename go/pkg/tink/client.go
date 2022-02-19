package tink

import (
	"net/http"
)

const (
	host              = "https://api.tink.com"
	TinkActorClientID = "df05e4b379934cd09963197cc855bfe9"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
}
