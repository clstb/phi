package beanacount

import (
	"github.com/shopspring/decimal"
	"time"
)

type Metadata struct {
	Key   string
	Value string
}

type AccountType struct {
	FinancialInstitutionId string
	ID                     string
	Name                   string
}

type Open struct {
	Date     string
	Account  AccountType
	Metadata []Metadata
}

type Posting struct {
	Account   AccountType
	Units     AmountType
	Cost      AmountType
	PriceType string
	Price     AmountType
}

type AmountType struct {
	Decimal  decimal.Decimal
	Currency string
}

type Provider struct {
	FinancialInstitutionId string
	DisplayName            string
}

type Transaction struct {
	Date      time.Time
	Type      string
	Payee     string
	Narration string
	Postings  []Posting
	Metadata  []Metadata
}

type Transactions []Transaction

type Value struct {
	Scale         int32
	UnscaledValue int64
}

type Amount struct {
	CurrencyCode string
	Value        Value
}

type Dates struct {
	Booked string
	Value  string
}

type TinkTransaction struct {
	Status       string
	AccountID    string
	ID           string
	Amount       Amount
	Dates        Dates
	Reference    string
	Descriptions string
}
