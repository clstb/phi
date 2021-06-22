package nordigen

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const host = "https://ob.nordigen.com/api"

type Client struct {
	client *http.Client
	token  string
}

func NewClient(token string) *Client {
	return &Client{
		token:  token,
		client: &http.Client{},
	}
}

type Bank struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	BIC       string   `json:"bic"`
	Countries []string `json:"countries"`
}

func (c *Client) GetBanks(country string) ([]Bank, error) {
	const endpoint = "/aspsps/"
	const method = http.MethodGet

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("country", country)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	banks := []Bank{}
	if err := json.NewDecoder(res.Body).Decode(&banks); err != nil {
		return nil, err
	}

	return banks, nil
}
