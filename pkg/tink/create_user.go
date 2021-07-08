package tink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type CreateUserReq struct {
	ExternalUserID string `json:"external_user_id"`
	Market         string `json:"market"`
	Locale         string `json:"locale"`
}

type CreateUserRes struct {
	ExternalUserID string `json:"external_user_id"`
	UserID         string `json:"user_id"`
}

var ErrUserExists = fmt.Errorf("user exists")

func (c *Client) CreateUser(body *CreateUserReq) (*CreateUserRes, error) {
	const endpoint = "/api/v1/user/create"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(body); err != nil {
		return nil, err
	}

	httpRes, err := c.client.Post(u.String(), "application/json", b)
	if err != nil {
		return nil, err
	}

	switch httpRes.StatusCode {
	case http.StatusOK:
	case http.StatusConflict:
		return nil, fmt.Errorf("status %d != 200: %w", httpRes.StatusCode, ErrUserExists)
	default:
		return nil, fmt.Errorf("status %d != 200: unknown error", httpRes.StatusCode)
	}
	if httpRes.StatusCode != http.StatusOK {
	}

	res := &CreateUserRes{}
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}
