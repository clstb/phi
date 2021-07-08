package tink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OAuthTokenReq struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}

type OAuthTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (c *Client) OAuthToken(body *OAuthTokenReq) (*OAuthTokenRes, error) {
	const endpoint = "/api/v1/oauth/token"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.client.PostForm(u.String(), url.Values{
		"code":          {body.Code},
		"client_id":     {body.ClientID},
		"client_secret": {body.ClientSecret},
		"grant_type":    {body.GrantType},
		"scope":         {body.Scope},
	})
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(httpRes.Body)
		fmt.Println(string(body))
		return nil, fmt.Errorf("status %d != 200", httpRes.StatusCode)
	}

	res := &OAuthTokenRes{}
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}
