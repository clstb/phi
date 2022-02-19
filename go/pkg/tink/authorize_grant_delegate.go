package tink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AuthorizeGrantDelegateReq struct {
	ResponseType   string
	ActorClientID  string
	UserID         string
	ExternalUserID string
	IDHint         string
	Scope          string
}

func (c *Client) AuthorizeGrantDelegate(body *AuthorizeGrantDelegateReq) (string, error) {
	const endpoint = "/api/v1/oauth/authorization-grant/delegate"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return "", err
	}

	httpRes, err := c.client.PostForm(u.String(), url.Values{
		"response_type":    {body.ResponseType},
		"actor_client_id":  {body.ActorClientID},
		"user_id":          {body.UserID},
		"external_user_id": {body.ExternalUserID},
		"id_hint":          {body.IDHint},
		"scope":            {body.Scope},
	})
	if err != nil {
		return "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(httpRes.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("status %d != 200: %s", httpRes.StatusCode, string(b))
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
