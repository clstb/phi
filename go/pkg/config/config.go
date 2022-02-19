package config

type Token struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresAt    int64
	Scope        string
}

type PhiToken Token
type TinkToken Token

type Config struct {
	PhiToken   PhiToken
	TinkToken  TinkToken
	Identities []string
	Recipients []string
}
