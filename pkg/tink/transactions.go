package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Transaction struct {
	ID        string `json:"id"`
	AccountID string `json:"accountId"`
	Amount    struct {
		Value struct {
			UnscaledValue int64 `json:"unscaledValue,string"`
			Scale         int32 `json:"scale,string"`
		} `json:"value"`
		CurrencyCode string `json:"currencyCode"`
	} `json:"amount"`
	Descriptions struct {
		Original           string `json:"original"`
		DisplayDescription string `json:"displayDescription"`
	} `json:"descriptions"`
	Dates struct {
		Booked string `json:"booked"`
	} `json:"dates"`
	Types struct {
		Type string `json:"type"`
	} `json:"types"`
	Categories struct {
		Pfm struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"pfm"`
	} `json:"categories"`
	Status string `json:"status"`
}

type TransactionsRes struct {
	Transactions  []Transaction `json:"transactions"`
	NextPageToken string        `json:"nextPageToken"`
}

func (c *Client) Transactions(
	token string,
	pageToken string,
) (*TransactionsRes, error) {
	const endpoint = "/data/v2/transactions"
	const method = http.MethodGet

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	if pageToken != "" {
		q.Add("pageToken", pageToken)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d != 200", res.StatusCode)
	}

	transactions := &TransactionsRes{}

	if err := json.NewDecoder(res.Body).Decode(transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
