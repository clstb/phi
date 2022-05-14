package client

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/tinkgw/pkg/client/rt"
	"github.com/clstb/phi/go/tinkgw/pkg/client/tink"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

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

func NewClient(url string) *Client {
	httpClient := &http.Client{Transport: transport}

	tinkClient := tink.NewClient(url+"/tink", tink.WithHTTPClient(httpClient))

	oryConf := ory.NewConfiguration()
	oryConf.Debug = true
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

type Client struct {
	*http.Client
	ctx        context.Context
	url        string
	tinkClient *tink.Client
	oryClient  *ory.APIClient
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Printf("Error: %s\n\n", err)
	}
}

func (c *Client) SendRequest(method string, url string, contentType string, body io.Reader) (string, error) {
	fmt.Println("Sending Request ------------------->")
	req, err := http.NewRequest(method, url, body)
	debug(httputil.DumpRequestOut(req, true))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	response, err := c.Do(req)
	if err == nil {
		fmt.Println("Received Response <------------------")
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		debug(body, err)
		return string(body), nil
	}
	return "", err
}
