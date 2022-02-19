package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OAuthTokenReq struct {
	Code         string `json:"code"`
	RefreshToken string `json:"refresh_token"`
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

	form := url.Values{
		"client_id":     {body.ClientID},
		"client_secret": {body.ClientSecret},
		"grant_type":    {body.GrantType},
		"scope":         {body.Scope},
	}
	switch body.GrantType {
	case "refresh_token":
		form.Add("refresh_token", body.RefreshToken)
	case "authorization_code":
		form.Add("code", body.Code)
	}

	res, err := c.client.PostForm(u.String(), form)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		tokenRes := &OAuthTokenRes{}
		if err := json.NewDecoder(res.Body).Decode(tokenRes); err != nil {
			return nil, err
		}
		return tokenRes, nil
	case http.StatusNotFound:
		return nil, UserNotFoundErr
	default:
		return nil, fmt.Errorf("unhandled status: %d", res.StatusCode)
	}
}
