package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Account struct {
	Balances struct {
		Booked struct {
			Amount struct {
				CurrencyCode string `json:"currencyCode"`
				Value        struct {
					Scale         int32 `json:"scale,string"`
					UnscaledValue int64 `json:"unscaledValue,string"`
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

func (c *Client) GetAccounts(pageToken string) ([]Account, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/data/v2/accounts?pageToken="+pageToken, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unhandled status: %d", res.StatusCode)
	}

	type Accounts struct {
		Accounts      []Account `json:"accounts"`
		NextPageToken string    `json:"nextPageToken"`
	}
	accounts := Accounts{}
	if err := json.NewDecoder(res.Body).Decode(&accounts); err != nil {
		return nil, err
	}

	if accounts.NextPageToken != "" {
		nextAccounts, err := c.GetAccounts(accounts.NextPageToken)
		if err != nil {
			return nil, err
		}
		accounts.Accounts = append(accounts.Accounts, nextAccounts...)
	}

	return accounts.Accounts, nil
}
