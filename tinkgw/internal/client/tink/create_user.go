package tink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clstb/phi/tinkgw/internal/config"
	"net/http"
)

var ErrUserExists = fmt.Errorf("user exists")

type CreateUserResponse struct {
	ExternalUserID string `json:"external_user_id"`
	UserID         string `json:"user_id"`
}

func (c *Client) CreateUserWithDefaults() (user CreateUserResponse, err error) {
	return c.CreateUser(config.DefaultMarket, config.DefaultLocale)
}

func (c *Client) CreateUser(market, locale string) (user CreateUserResponse, err error) {
	b := &bytes.Buffer{}
	if err = json.NewEncoder(b).Encode(map[string]string{
		"market": market,
		"locale": locale,
	}); err != nil {
		return
	}

	res, err := c.Post(config.TinkApiUri+config.UserCreatePath, config.JsonMediaType, b)
	if err != nil {
		return user, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(res.Body).Decode(&user)
	case http.StatusConflict:
		err = ErrUserExists
	default:
		err = fmt.Errorf("unimplemented status: %d", res.StatusCode)
	}

	return
}
