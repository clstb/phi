package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetAuthorizeGrantDelegateCode(
	responseType,
	userId,
	externalUserId,
	idHint,
	scope string,
) (s string, err error) {
	res, err := c.PostForm(c.url+"/api/v1/oauth/authorization-grant/delegate", url.Values{
		"response_type":    {responseType},
		"actor_client_id":  {"df05e4b379934cd09963197cc855bfe9"},
		"user_id":          {userId},
		"external_user_id": {externalUserId},
		"id_hint":          {idHint},
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
