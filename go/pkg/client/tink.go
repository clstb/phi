package client

import "github.com/clstb/phi/go/pkg/client/tink"

func (c *Client) GetAccounts(pageToken string) ([]tink.Account, error) {
	return c.tinkClient.GetAccounts(pageToken)
}

func (c *Client) GetTransactions(pageToken string) ([]tink.Transaction, error) {
	return c.tinkClient.GetTransactions(pageToken)
}

func (c *Client) GetProviders(country string) ([]tink.Provider, error) {
	return c.tinkClient.GetProviders(country)
}
