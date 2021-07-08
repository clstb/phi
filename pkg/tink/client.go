package tink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	host              = "https://api.tink.com"
	TinkActorClientID = "df05e4b379934cd09963197cc855bfe9"
)

type Client struct {
	client *http.Client
}

func NewClient(
	clientID string,
	clientSecret string,
	scope string,
) (*Client, error) {
	client := &http.Client{}
	c := &Client{
		client: client,
	}

	token, err := c.OAuthToken(&OAuthTokenReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
		Scope:        scope,
	})
	if err != nil {
		return nil, err
	}

	c.client.Transport = &AuthorizationRoundTripper{
		Token: token.AccessToken,
		Next:  http.DefaultTransport,
	}

	return c, nil
}

type OAuthTokenReq struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}

type OAuthTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (c *Client) OAuthToken(body *OAuthTokenReq) (*OAuthTokenRes, error) {
	const endpoint = "/api/v1/oauth/token"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.client.PostForm(u.String(), url.Values{
		"code":          {body.Code},
		"client_id":     {body.ClientID},
		"client_secret": {body.ClientSecret},
		"grant_type":    {body.GrantType},
		"scope":         {body.Scope},
	})
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(httpRes.Body)
		fmt.Println(string(body))
		return nil, fmt.Errorf("status %d != 200", httpRes.StatusCode)
	}

	res := &OAuthTokenRes{}
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

type AuthorizeGrantDelegateReq struct {
	ResponseType   string
	ActorClientID  string
	UserID         string
	ExternalUserID string
	IDHint         string
	Scope          string
}

func (c *Client) AuthorizeGrantDelegate(body *AuthorizeGrantDelegateReq) (string, error) {
	const endpoint = "/api/v1/oauth/authorization-grant/delegate"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return "", err
	}

	httpRes, err := c.client.PostForm(u.String(), url.Values{
		"response_type":    {body.ResponseType},
		"actor_client_id":  {body.ActorClientID},
		"user_id":          {body.UserID},
		"external_user_id": {body.ExternalUserID},
		"id_hint":          {body.IDHint},
		"scope":            {body.Scope},
	})
	if err != nil {
		return "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d != 200", httpRes.StatusCode)
	}

	type response struct {
		Code string `json:"code"`
	}
	res := &response{}
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return "", nil
	}

	return res.Code, nil
}

type AuthorizeGrantReq struct {
	UserID         string `json:"user_id"`
	ExternalUserID string `json:"external_user_id"`
	Scope          string `json:"scope"`
}

func (c *Client) AuthorizeGrant(body *AuthorizeGrantReq) (string, error) {
	const endpoint = "/api/v1/oauth/authorization-grant"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return "", err
	}

	httpRes, err := c.client.PostForm(u.String(), url.Values{
		"user_id":          {body.UserID},
		"external_user_id": {body.ExternalUserID},
		"scope":            {body.Scope},
	})
	if err != nil {
		return "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d != 200", httpRes.StatusCode)
	}

	type response struct {
		Code string `json:"code"`
	}
	res := &response{}
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return "", nil
	}

	return res.Code, nil
}

type CreateUserReq struct {
	ExternalUserID string `json:"external_user_id"`
	Market         string `json:"market"`
	Locale         string `json:"locale"`
}

type CreateUserRes struct {
	ExternalUserID string `json:"external_user_id"`
	UserID         string `json:"user_id"`
}

var ErrUserExists = fmt.Errorf("user exists")

func (c *Client) CreateUser(body *CreateUserReq) (*CreateUserRes, error) {
	const endpoint = "/api/v1/user/create"
	const method = http.MethodPost

	u, err := url.Parse(host + endpoint)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(body); err != nil {
		return nil, err
	}

	httpRes, err := c.client.Post(u.String(), "application/json", b)
	if err != nil {
		return nil, err
	}

	switch httpRes.StatusCode {
	case http.StatusOK:
	case http.StatusConflict:
		return nil, fmt.Errorf("status %d != 200: %w", httpRes.StatusCode, ErrUserExists)
	default:
		return nil, fmt.Errorf("status %d != 200: unknown error", httpRes.StatusCode)
	}
	if httpRes.StatusCode != http.StatusOK {
	}

	res := &CreateUserRes{}
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

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
