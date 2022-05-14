package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id string `json:"id"`
}

func (c *Client) GetUser() (user User, err error) {
	res, err := c.httpClient.Get(c.url + "/api/v1/user")
	if err != nil {
		return user, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(res.Body).Decode(&user)
	default:
		err = fmt.Errorf("unimplemented status: %d", res.StatusCode)
	}

	return
}
