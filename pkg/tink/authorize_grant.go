package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type AuthorizeGrantReq struct {
	UserID         string `json:"user_id"`
	ExternalUserID string `json:"external_user_id"`
	Scope          string `json:"scope"`
}

func (c *Client) AuthorizeGrant(body *AuthorizeGrantReq) (string, error) {
	const endpoint = "/api/v1/oauth/authorization-grant"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return "", err
	}

	httpRes, err := c.client.PostForm(u.String(), url.Values{
		"user_id":          {body.UserID},
		"external_user_id": {body.ExternalUserID},
		"scope":            {body.Scope},
	})
	if err != nil {
		return "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d != 200", httpRes.StatusCode)
	}

	type response struct {
		Code string `json:"code"`
	}
	res := &response{}
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return "", nil
	}

	return res.Code, nil
}
