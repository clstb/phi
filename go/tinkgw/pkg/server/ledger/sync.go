package ledger

import (
	"fmt"
	"github.com/clstb/phi/go/tinkgw/pkg/client/tink"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

func (l *Ledger) UpdateLedger(providers []tink.Provider, accounts []tink.Account, transactions []tink.Transaction) {
	providersById := map[string]tink.Provider{}
	for _, provider := range providers {
		providersById[provider.FinancialInstitutionID] = provider
	}

	opensByTinkId := l.Opens().ByTinkId()
	for _, account := range accounts {
		_, ok := opensByTinkId[account.ID]
		if ok {
			continue
		}

		*l = append(*l, Open{
			Date: "1970-01-01",
			Account: fmt.Sprintf(
				"Assets:%s:%s",
				providersById[account.FinancialInstitutionID].DisplayName,
				account.Name,
			),
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

		amount := Amount{
			Decimal: decimal.New(
				transaction.Amount.Value.UnscaledValue,
				transaction.Amount.Value.Scale*-1,
			),
			Currency: transaction.Amount.CurrencyCode,
		}
		var balanceAccount string
		if amount.IsNegative() {
			balanceAccount = "Expenses:Unassigned"
		} else {
			balanceAccount = "Income:Unassigned"
		}

		date, _ := time.Parse("2006-01-02", transaction.Dates.Booked)

		*l = append(*l, Transaction{
			Date:      date,
			Type:      "*",
			Payee:     transaction.Reference,
			Narration: transaction.Descriptions.Display,
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
