package tink

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/tinkgw/pkg/client/rt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	*http.Client
	ctx context.Context
	url string
}

func (c *Client) GetLink() (s string, err error) {
	url, err := url.Parse(c.url + "/link")
	if err != nil {
		return
	}

	httpResp, err := c.Get(url.String())
	if err != nil {
		return
	}

	switch httpResp.StatusCode {
	case http.StatusOK:
		var b []byte
		b, err = ioutil.ReadAll(httpResp.Body)
		s = string(b)
	default:
		err = fmt.Errorf("unhandled status: %d", httpResp.StatusCode)
	}

	return
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

func NewClient(url string) *Client {
	httpClient := &http.Client{Transport: transport}

	return &Client{
		Client: httpClient,
		ctx:    context.Background(),
		url:    url,
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
