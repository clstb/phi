package tink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ProviderConsentsRes struct {
	ProviderConsents []struct {
		AccountIds    []string `json:"accountIds"`
		CredentialsID string   `json:"credentialsId"`
		DetailedError struct {
			Details struct {
				Reason    string `json:"reason"`
				Retryable bool   `json:"retryable"`
			} `json:"details"`
			DisplayMessage string `json:"displayMessage"`
			Type           string `json:"type"`
		} `json:"detailedError"`
		ProviderName      string `json:"providerName"`
		SessionExpiryDate int64  `json:"sessionExpiryDate"`
		Status            string `json:"status"`
		StatusUpdated     int64  `json:"statusUpdated"`
	} `json:"providerConsents"`
}

func (c *Client) ProviderConsents(
	token string,
) (*ProviderConsentsRes, error) {
	const endpoint = "/api/v1/provider-consents"
	const method = http.MethodGet

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

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

	providerConsents := &ProviderConsentsRes{}

	if err := json.NewDecoder(res.Body).Decode(providerConsents); err != nil {
		return nil, err
	}

	return providerConsents, nil
}
