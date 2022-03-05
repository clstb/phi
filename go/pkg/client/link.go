package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

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
