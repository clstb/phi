package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Account struct {
	Balances struct {
		Booked struct {
			Amount struct {
				CurrencyCode string `json:"currencyCode"`
				Value        struct {
					Scale         string `json:"scale"`
					UnscaledValue string `json:"unscaledValue"`
				} `json:"value"`
			} `json:"amount"`
		} `json:"booked"`
	} `json:"balances"`
	CustomerSegment string `json:"customerSegment"`
	Dates           struct {
		LastRefreshed time.Time `json:"lastRefreshed"`
	} `json:"dates"`
	FinancialInstitutionID string `json:"financialInstitutionId"`
	ID                     string `json:"id"`
	Identifiers            struct {
		FinancialInstitution struct {
			AccountNumber string `json:"accountNumber"`
		} `json:"financialInstitution"`
		Iban struct {
			Bban string `json:"bban"`
			Iban string `json:"iban"`
		} `json:"iban"`
		Pan struct {
			Masked string `json:"masked"`
		} `json:"pan"`
	} `json:"identifiers"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type AccountsRes struct {
	Accounts      []Account `json:"accounts"`
	NextPageToken string    `json:"nextPageToken"`
}

func (c *Client) Accounts(
	token string,
	pageToken string,
) (*AccountsRes, error) {
	const endpoint = "/data/v2/accounts"
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

	accounts := &AccountsRes{}

	if err := json.NewDecoder(res.Body).Decode(accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}
