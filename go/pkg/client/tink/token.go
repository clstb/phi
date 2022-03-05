package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (c *Client) GetToken(
	code,
	refreshToken,
	clientId,
	clientSecret,
	grantType,
	scope string,
) (token Token, err error) {
	form := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"grant_type":    {grantType},
		"scope":         {scope},
	}
	switch grantType {
	case "refresh_token":
		form.Add("refresh_token", refreshToken)
	case "authorization_code":
		form.Add("code", code)
	}

	res, err := c.PostForm(c.url+"/api/v1/oauth/token", form)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(res.Body).Decode(&token)
	default:
		err = fmt.Errorf("unimplemented status: %d", res.StatusCode)
	}

	return
}
