package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) GetTransactions(pageToken string) ([]Transaction, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/..data/v2/transactions?pageToken="+pageToken, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unhandled status: %d", res.StatusCode)
	}

	type Transactions struct {
		Transactions  []Transaction `json:"transactions"`
		NextPageToken string        `json:"nextPageToken"`
	}
	transactions := Transactions{}
	if err := json.NewDecoder(res.Body).Decode(&transactions); err != nil {
		return nil, err
	}

	if transactions.NextPageToken != "" {
		nextTransactions, err := c.GetTransactions(transactions.NextPageToken)
		if err != nil {
			return nil, err
		}
		transactions.Transactions = append(transactions.Transactions, nextTransactions...)
	}

	return transactions.Transactions, nil
}
