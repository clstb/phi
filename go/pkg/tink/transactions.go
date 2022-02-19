package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Transaction struct {
	AccountID string `json:"accountId"`
	Amount    struct {
		CurrencyCode string `json:"currencyCode"`
		Value        struct {
			Scale         int32 `json:"scale,string"`
			UnscaledValue int64 `json:"unscaledValue,string"`
		} `json:"value"`
	} `json:"amount"`
	Categories struct {
		Pfm struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"pfm"`
	} `json:"categories"`
	Dates struct {
		Booked string `json:"booked"`
		Value  string `json:"value"`
	} `json:"dates"`
	Descriptions struct {
		Display  string `json:"display"`
		Original string `json:"original"`
	} `json:"descriptions"`
	ID          string `json:"id"`
	Identifiers struct {
		ProviderTransactionID string `json:"providerTransactionId"`
	} `json:"identifiers"`
	MerchantInformation struct {
		MerchantCategoryCode string `json:"merchantCategoryCode"`
		MerchantName         string `json:"merchantName"`
	} `json:"merchantInformation"`
	ProviderMutability string `json:"providerMutability"`
	Reference          string `json:"reference"`
	Status             string `json:"status"`
	Types              struct {
		FinancialInstitutionTypeCode string `json:"financialInstitutionTypeCode"`
		Type                         string `json:"type"`
	} `json:"types"`
}

type Transactions []Transaction

type TransactionsRes struct {
	Transactions  Transactions `json:"transactions"`
	NextPageToken string       `json:"nextPageToken"`
}

func (c *Client) Transactions(pageToken string) (*TransactionsRes, error) {
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

	res, err := c.client.Do(req)
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
