package nordigen

import (
	"bytes"
	"encoding/json"
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
	req.Header.Add("Authorization", "Token "+c.token)

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

type EndUserAgreement struct {
	ID                 string `json:"id"`
	Created            string `json:"created"`
	Accepted           bool   `json:"accepted"`
	MaxHistoricalDays  int    `json:"max_historical_days"`
	AccessValidForDays int    `json:"access_valid_for_days"`
	EndUserID          string `json:"enduser_id"`
	AspspID            string `json:"aspsp_id"`
}

func (c *Client) CreateEndUserAgreement(
	maxHistoricalDays string,
	endUserID string,
	aspspID string,
) (*EndUserAgreement, error) {
	const endpoint = "/agreements/enduser/"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

	reqBody := struct {
		MaxHistoricalDays string `json:"max_historical_days"`
		EndUserID         string `json:"enduser_id"`
		AspspID           string `json:"aspsp_id"`
	}{
		MaxHistoricalDays: maxHistoricalDays,
		EndUserID:         endUserID,
		AspspID:           aspspID,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(reqBody); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), &b)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+c.token)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	eua := &EndUserAgreement{}
	if err := json.NewDecoder(res.Body).Decode(eua); err != nil {
		return nil, err
	}

	return eua, nil
}
