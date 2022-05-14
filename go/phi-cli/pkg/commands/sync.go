package commands

import (
	"fmt"
	beancount2 "github.com/clstb/phi/go/phi-cli/pkg/beancount"
	"github.com/clstb/phi/go/tinkgw/pkg/client"
	"github.com/clstb/phi/go/tinkgw/pkg/client/tink"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shopspring/decimal"
)

func Sync(
	ledger beancount2.Ledger,
	client *client.Client,
) tea.Cmd {
	return func() tea.Msg {
		providers, err := client.GetProviders("DE")
		if err != nil {
			return err
		}

		accounts, err := client.GetAccounts("")
		if err != nil {
			return err
		}

		transactions, err := client.GetTransactions("")
		if err != nil {
			return err
		}

		var filteredTransactions []tink.Transaction
		for _, transaction := range transactions {
			if transaction.Status != "BOOKED" {
				continue
			}
			filteredTransactions = append(filteredTransactions, transaction)
		}

		return UpdateLedger(ledger, providers, accounts, filteredTransactions)
	}
}

func UpdateLedger(
	ledger beancount2.Ledger,
	providers []tink.Provider,
	accounts []tink.Account,
	transactions []tink.Transaction,
) beancount2.Ledger {
	providersById := map[string]tink.Provider{}
	for _, provider := range providers {
		providersById[provider.FinancialInstitutionID] = provider
	}

	opensByTinkId := ledger.Opens().ByTinkId()
	for _, account := range accounts {
		_, ok := opensByTinkId[account.ID]
		if ok {
			continue
		}

		ledger = append(ledger, beancount2.Open{
			Date: "1970-01-01",
			Account: fmt.Sprintf(
				"Assets:%s:%s",
				providersById[account.FinancialInstitutionID].DisplayName,
				account.Name,
			),
			Metadata: []beancount2.Metadata{
				{
					Key:   "tink_id",
					Value: strconv.Quote(account.ID),
				},
			},
		})
	}
	opensByTinkId = ledger.Opens().ByTinkId()

	transactionsByTinkId := ledger.Transactions().ByTinkId()
	for _, transaction := range transactions {
		_, ok := transactionsByTinkId[transaction.ID]
		if ok {
			continue
		}

		amount := beancount2.Amount{
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

		ledger = append(ledger, beancount2.Transaction{
			Date:      date,
			Type:      "*",
			Payee:     transaction.Reference,
			Narration: transaction.Descriptions.Display,
			Metadata: []beancount2.Metadata{
				{
					Key:   "tink_id",
					Value: strconv.Quote(transaction.ID),
				},
			},
			Postings: []beancount2.Posting{
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
	return ledger
}
