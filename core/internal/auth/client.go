package auth

import (
	"context"
	"github.com/clstb/phi/core/internal/config"
	ory "github.com/ory/kratos-client-go"
	"net"
	"net/http"
	"time"
)

type AuthClient struct {
	*http.Client
	ctx       context.Context
	url       string
	OryClient *ory.APIClient
}

type Session struct {
	ory.Session
	Token string
}

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

func NewClient(url string) *AuthClient {
	httpClient := &http.Client{Transport: transport}

	oryConf := ory.NewConfiguration()
	oryConf.Debug = true
	oryConf.HTTPClient = httpClient
	oryConf.Servers = []ory.ServerConfiguration{{URL: url + config.OryPath}}
	oryClient := ory.NewAPIClient(oryConf)

	return &AuthClient{
		Client:    httpClient,
		ctx:       context.Background(),
		url:       url,
		OryClient: oryClient,
	}
}

func (c *AuthClient) URL() string {
	return c.url
}
