package tink

import (
	"fmt"
	"github.com/clstb/phi/pkg"
	"github.com/jobala/middleware_pipeline/pipeline"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type Client struct {
	*http.Client
	logger *zap.SugaredLogger
}
type AuthorizationMiddleware struct {
	bearerToken string
}

func (s AuthorizationMiddleware) Intercept(pipeline pipeline.Pipeline, req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+s.bearerToken)
	body, _ := httputil.DumpRequest(req, true)
	log.Println(fmt.Sprintf("%s", string(body)))
	return pipeline.Next(req)
}

func NewAuthorizedClient(token string, logger *zap.SugaredLogger) *Client {
	transport := pipeline.NewCustomTransport(&AuthorizationMiddleware{token}, &pkg.LoggingMiddleware{})
	transport.ForceAttemptHTTP2 = true
	transport.MaxIdleConns = 10
	transport.IdleConnTimeout = 30 * time.Second
	transport.IdleConnTimeout = 90 * time.Second

	httpClient := &http.Client{Transport: transport}
	return &Client{
		Client: httpClient,
		logger: logger.Named("tink-client-authorized"),
	}
}

func NewClient(logger *zap.SugaredLogger) *Client {
	transport := pipeline.NewCustomTransport(&pkg.LoggingMiddleware{})
	transport.ForceAttemptHTTP2 = true
	transport.MaxIdleConns = 10
	transport.IdleConnTimeout = 30 * time.Second
	transport.IdleConnTimeout = 90 * time.Second

	httpClient := &http.Client{Transport: transport}
	return &Client{
		Client: httpClient,
		logger: logger.Named("tink-client-basic"),
	}
}
