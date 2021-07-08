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

func NewClient(
	clientID string,
	clientSecret string,
	scope string,
) (*Client, error) {
	client := &http.Client{}
	c := &Client{
		client: client,
	}

	token, err := c.OAuthToken(&OAuthTokenReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
		Scope:        scope,
	})
	if err != nil {
		return nil, err
	}

	c.client.Transport = &AuthorizationRoundTripper{
		Token: token.AccessToken,
		Next:  http.DefaultTransport,
	}

	return c, nil
}
