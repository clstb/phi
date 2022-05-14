package tink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clstb/phi/go/tinkgw/pkg"
	"net/http"
)

type CreateUserResponse struct {
	ExternalUserID string `json:"external_user_id"`
	UserID         string `json:"user_id"`
}

func (c *Client) CreateUser(
	externalUserId,
	market,
	locale string,
) (user CreateUserResponse, err error) {
	b := &bytes.Buffer{}
	if err = json.NewEncoder(b).Encode(map[string]string{
		"external_user_id": externalUserId,
		"market":           market,
		"locale":           locale,
	}); err != nil {
		return
	}

	res, err := c.Post(c.url+"/api/v1/user/create", "application/json", b)
	if err != nil {
		pkg.Sugar.Error(err)
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
