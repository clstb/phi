package tink

import (
	"net"
	"net/http"
	"time"

	"github.com/clstb/phi/go/pkg/client/rt"
)

var transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	DisableCompression:    true,
}

type Client struct {
	*http.Client
	url string
}

type Opt func(*Client)

func WithHTTPClient(httpClient *http.Client) func(*Client) {
	return func(c *Client) {
		c.Client = httpClient
	}
}

func NewClient(url string, opts ...Opt) *Client {
	httpClient := &http.Client{Transport: transport}

	client := &Client{
		Client: httpClient,
		url:    url,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) SetBearerToken(token string) {
	c.Transport = rt.AuthorizationRoundTripper{
		Token: token,
		Next:  transport,
	}
}
