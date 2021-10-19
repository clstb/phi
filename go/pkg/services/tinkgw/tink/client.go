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

type Opt func(*Client) error

func WithAuth(
	clientID,
	clientSecret,
	scope string,
) func(*Client) error {
	return func(client *Client) error {
		token, err := client.OAuthToken(&OAuthTokenReq{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			GrantType:    "client_credentials",
			Scope:        scope,
		})
		if err != nil {
			return err
		}

		client.client.Transport = &AuthorizationRoundTripper{
			Token: token.AccessToken,
			Next:  http.DefaultTransport,
		}

		return nil
	}
}

func NewClient(opts ...Opt) (*Client, error) {
	client := &Client{
		client: &http.Client{},
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}
