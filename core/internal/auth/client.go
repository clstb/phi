package auth

import (
	"context"
	"github.com/clstb/phi/pkg"
	"github.com/jobala/middleware_pipeline/pipeline"
	ory "github.com/ory/kratos-client-go"
	"net/http"
	"time"
)

type Client struct {
	*http.Client
	ctx       context.Context
	OryClient *ory.APIClient
}

type Session struct {
	ory.Session
	Token string
}

func NewClient(oryUri string) *Client {

	transport := pipeline.NewCustomTransport(&pkg.LoggingMiddleware{})
	transport.ForceAttemptHTTP2 = true
	transport.MaxIdleConns = 10
	transport.IdleConnTimeout = 30 * time.Second
	transport.IdleConnTimeout = 90 * time.Second

	httpClient := &http.Client{Transport: transport}

	oryConf := ory.NewConfiguration()
	oryConf.Debug = true
	oryConf.HTTPClient = httpClient
	oryConf.Servers = []ory.ServerConfiguration{{URL: oryUri}}
	oryClient := ory.NewAPIClient(oryConf)

	return &Client{
		Client:    httpClient,
		ctx:       context.Background(),
		OryClient: oryClient,
	}
}
