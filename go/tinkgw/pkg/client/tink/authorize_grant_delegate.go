package tink

import (
	"encoding/json"
	"fmt"
	"github.com/clstb/phi/go/tinkgw/pkg/config"
	"net/http"
	"net/url"
)

type Code struct {
	Code string `json:"code"`
}

func (c *Client) GetDelegatedAutorizationCode(actorClientId string, userId string) (s string, err error) {
	res, err := c.PostForm(c.url+config.DelegatedAuthorizationEndpoint, url.Values{
		"actor_client_id": {actorClientId},
		"user_id":         {userId},
		"id_hint":         {config.GetAuthorizeGrantDelegateCodeRoles},
		"scope":           {config.GetAuthorizeGrantDelegateCodeRoles},
	})
	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusOK {
		code := Code{}
		s, err = code.Code, json.NewDecoder(res.Body).Decode(&code)
		return s, err
	}
	return "", fmt.Errorf("unimplemented status: %d", res.StatusCode)

}
