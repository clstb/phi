package tink

import (
	"encoding/json"
	"fmt"
	"github.com/clstb/phi/tinkgw/internal/config"
	"io/ioutil"
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

type Accounts struct {
	Accounts      []Account `json:"accounts"`
	NextPageToken string    `json:"nextPageToken"`
}

func (c *Client) GetAccounts() ([]Account, error) {
	res, err := c.Get(config.TinkApiUri + config.AccountsPath)

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	byteArr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var accounts Accounts
	err = json.Unmarshal(byteArr, &accounts)
	fmt.Println(string(byteArr))
	if err != nil {
		return nil, err
	}

	return accounts.Accounts, nil

	/*
		if accounts.NextPageToken != "" {
			nextAccounts, err := c.GetAccounts(accounts.NextPageToken)
			if err != nil {
				return nil, err
			}
			accounts.Accounts = append(accounts.Accounts, nextAccounts...)
		}*/
}
