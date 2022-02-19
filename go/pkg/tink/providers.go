package tink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Provider struct {
	AccessType             string   `json:"accessType"`
	AuthenticationFlow     string   `json:"authenticationFlow"`
	AuthenticationUserType string   `json:"authenticationUserType"`
	Capabilities           []string `json:"capabilities"`
	CredentialsType        string   `json:"credentialsType"`
	Currency               string   `json:"currency"`
	DisplayDescription     string   `json:"displayDescription"`
	DisplayName            string   `json:"displayName"`
	Fields                 []struct {
		AdditionalInfo string   `json:"additionalInfo"`
		Checkbox       bool     `json:"checkbox"`
		DefaultValue   string   `json:"defaultValue"`
		Description    string   `json:"description"`
		Group          string   `json:"group"`
		HelpText       string   `json:"helpText"`
		Hint           string   `json:"hint"`
		Immutable      bool     `json:"immutable"`
		Masked         bool     `json:"masked"`
		MaxLength      int      `json:"maxLength"`
		MinLength      int      `json:"minLength"`
		Name           string   `json:"name"`
		Numeric        bool     `json:"numeric"`
		OneOf          bool     `json:"oneOf"`
		Optional       bool     `json:"optional"`
		Options        []string `json:"options"`
		Pattern        string   `json:"pattern"`
		PatternError   string   `json:"patternError"`
		SelectOptions  []struct {
			IconURL string `json:"iconUrl"`
			Text    string `json:"text"`
			Value   string `json:"value"`
		} `json:"selectOptions"`
		Sensitive bool   `json:"sensitive"`
		Style     string `json:"style"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"fields"`
	FinancialInstitutionID   string `json:"financialInstitutionId"`
	FinancialInstitutionName string `json:"financialInstitutionName"`
	FinancialServices        []struct {
		Segment   string `json:"segment"`
		ShortName string `json:"shortName"`
	} `json:"financialServices"`
	GroupDisplayName string `json:"groupDisplayName"`
	Images           struct {
		Banner string `json:"banner"`
		Icon   string `json:"icon"`
	} `json:"images"`
	LoginHeaderColour string   `json:"loginHeaderColour"`
	Market            string   `json:"market"`
	MultiFactor       bool     `json:"multiFactor"`
	Name              string   `json:"name"`
	PasswordHelpText  string   `json:"passwordHelpText"`
	PisCapabilities   []string `json:"pisCapabilities"`
	Popular           bool     `json:"popular"`
	ReleaseStatus     string   `json:"releaseStatus"`
	Status            string   `json:"status"`
	Transactional     bool     `json:"transactional"`
	Type              string   `json:"type"`
}

type Providers []Provider

type ProvidersRes struct {
	Providers Providers `json:"providers"`
}

func (c *Client) Providers(country string) (*ProvidersRes, error) {
	const endpoint = "/api/v1/providers"
	const method = http.MethodGet

	u, err := url.Parse(host + endpoint + "/" + country)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("status %d != 200: %s", res.StatusCode, string(b))
	}

	providers := &ProvidersRes{}

	if err := json.NewDecoder(res.Body).Decode(providers); err != nil {
		return nil, err
	}

	return providers, nil
}
