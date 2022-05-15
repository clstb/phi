package tink

import (
	"encoding/json"
	"fmt"
	"github.com/clstb/phi/tinkgw/internal/config"
	"io/ioutil"
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

type Transactions struct {
	Transactions  []Transaction `json:"transactions"`
	NextPageToken string        `json:"nextPageToken"`
}

func (c *Client) GetTransactions() ([]Transaction, error) {

	res, err := c.httpClient.Get(c.url + config.TransactionsPath)

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

	var transaction Transactions
	err = json.Unmarshal(byteArr, &transaction)
	fmt.Println(string(byteArr))
	if err != nil {
		return nil, err
	}

	return transaction.Transactions, nil
}
