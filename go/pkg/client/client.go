package client

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/clstb/phi/go/pkg/client/rt"
	"github.com/clstb/phi/go/pkg/client/tink"
	ory "github.com/ory/kratos-client-go"
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
	ctx        context.Context
	url        string
	tinkClient *tink.Client
	oryClient  *ory.APIClient
}

func NewClient(url string) *Client {
	httpClient := &http.Client{Transport: transport}

	tinkClient := tink.NewClient(url+"/tink", tink.WithHTTPClient(httpClient))

	oryConf := ory.NewConfiguration()
	oryConf.HTTPClient = httpClient
	oryConf.Servers = []ory.ServerConfiguration{{URL: url + "/ory"}}
	oryClient := ory.NewAPIClient(oryConf)

	return &Client{
		Client:     httpClient,
		ctx:        context.Background(),
		url:        url,
		tinkClient: tinkClient,
		oryClient:  oryClient,
	}
}

func (c *Client) SetBearerToken(token string) {
	c.Transport = rt.AuthorizationRoundTripper{
		Token: token,
		Next:  transport,
	}
}

func (c *Client) URL() string {
	return c.url
}
