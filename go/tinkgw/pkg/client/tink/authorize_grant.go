package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetAuthorizeGrantCode(
	userId,
	externalUserId,
	scope string,
) (s string, err error) {
	res, err := c.httpClient.PostForm(c.url+"/api/v1/oauth/authorization-grant", url.Values{
		"user_id":          {userId},
		"external_user_id": {externalUserId},
		"scope":            {scope},
	})
	if err != nil {
		return "", err
	}

	switch res.StatusCode {
	case http.StatusOK:
		type Code struct {
			Code string `json:"code"`
		}
		code := Code{}
		s, err = code.Code, json.NewDecoder(res.Body).Decode(&code)
	default:
		err = fmt.Errorf("unimplemented status: %d", res.StatusCode)
	}

	return
}
