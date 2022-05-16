package beanacount

import (
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

func (l *Ledger) UpdateLedger(providers []Provider, accounts []AccountType, transactions []TinkTransaction) {
	providersById := map[string]string{}
	for _, provider := range providers {
		providersById[provider.FinancialInstitutionId] = provider.DisplayName
	}

	opensByTinkId := l.Opens().ByTinkId()
	for _, account := range accounts {
		_, ok := opensByTinkId[account.ID]
		if ok {
			continue
		}

		*l = append(*l, Open{
			Date: "1970-01-01",
			Account: AccountType{
				FinancialInstitutionId: "Assets",
				ID:                     providersById[account.FinancialInstitutionId],
				Name:                   account.Name,
			},
			Metadata: []Metadata{
				{
					Key:   "tink_id",
					Value: strconv.Quote(account.ID),
				},
			},
		})
	}
	opensByTinkId = l.Opens().ByTinkId()

	transactionsByTinkId := l.Transactions().ByTinkId()
	for _, transaction := range transactions {
		_, ok := transactionsByTinkId[transaction.ID]
		if ok {
			continue
		}

		amount := AmountType{
			Decimal: decimal.New(
				transaction.Amount.Value.UnscaledValue,
				transaction.Amount.Value.Scale*-1,
			),
			Currency: transaction.Amount.CurrencyCode,
		}
		var balanceAccount AccountType
		if amount.IsNegative() {
			balanceAccount = AccountType{FinancialInstitutionId: "Expenses", Name: "Unassigned"}
		} else {
			balanceAccount = AccountType{FinancialInstitutionId: "Income", Name: "Unassigned"}
		}

		date, _ := time.Parse("2006-01-02", transaction.Dates.Booked)

		*l = append(*l, Transaction{
			Date:      date,
			Type:      "*",
			Payee:     transaction.Reference,
			Narration: transaction.Descriptions,
			Metadata: []Metadata{
				{
					Key:   "tink_id",
					Value: strconv.Quote(transaction.ID),
				},
			},
			Postings: []Posting{
				{
					Account: balanceAccount,
					Units:   amount.Neg(),
				},
				{
					Account: opensByTinkId[transaction.AccountID].Account,
					Units:   amount,
				},
			},
		})
	}
}
