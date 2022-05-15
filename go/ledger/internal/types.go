package internal

import (
	"github.com/shopspring/decimal"
	"time"
)

type Metadata struct {
	Key   string
	Value string
}

type Open struct {
	Date     string
	Account  string
	Metadata []Metadata
}

type Posting struct {
	Account   string
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

type Account struct {
	FinancialInstitutionId string
	ID                     string
	Name                   string
}

type Transaction struct {
	Date      time.Time
	Type      string
	Payee     string
	Narration string
	Postings  []Posting
	Metadata  []Metadata
}

type TinkTransaction struct {
	Status    string
	AccountID string
	ID        string
	Amount    struct {
		CurrencyCode string
		Value        struct {
			Scale         int32
			UnscaledValue int64
		}
	}
	Dates struct {
		Booked string
		Value  string
	}
	Reference    string
	Descriptions string
}
