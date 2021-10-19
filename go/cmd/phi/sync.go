package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/clstb/phi/go/pkg/config"
	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/parser"
	"github.com/clstb/phi/go/pkg/services/tinkgw/pb"
	"github.com/clstb/phi/go/pkg/services/tinkgw/tink"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Sync(ctx *cli.Context) error {
	configPath := ctx.String("config")
	config, err := config.Load(configPath)
	if err != nil {
		return err
	}
	expiresAt := time.Unix(config.TinkToken.ExpiresAt, 0)

	if !time.Now().Before(expiresAt) {
		tinkGWHost := ctx.String("tinkgw-host")
		conn, err := grpc.Dial(
			tinkGWHost,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(interceptor.ClientAuthUnary(config.AccessToken)),
			grpc.WithStreamInterceptor(interceptor.ClientAuthStream(config.AccessToken)),
		)
		if err != nil {
			return err
		}

		token := &pb.Token{}
		if time.Now().Before(expiresAt.Add(time.Hour * 24)) {
			token.RefreshToken = config.TinkToken.RefreshToken
		}

		client := pb.NewTinkGWClient(conn)

		token, err = client.GetToken(ctx.Context, token)
		if err != nil {
			return err
		}
		config.TinkToken = token

		if err := config.Save(configPath); err != nil {
			return err
		}
	}

	return syncLedger(config.TinkToken.AccessToken)
}

func syncLedger(accessToken string) error {
	fetchedAccounts, err := fetchAccounts(accessToken)
	if err != nil {
		return err
	}

	fetchedTransactions, err := fetchTransactions(accessToken)
	if err != nil {
		return err
	}

	ledger, err := parser.Load("./bla.beancount")
	if err != nil {
		return err
	}

	accountsByTinkID := ledger.Opens().ByTinkID()
	transactionsByTinkID := ledger.Transactions().ByTinkID()

	for _, account := range fetchedAccounts {
		if _, ok := accountsByTinkID[fmt.Sprintf(`"%s"`, account.ID)]; ok {
			continue
		}

		ledger = append(ledger, parser.Open{
			Date: "1970-01-01",
			Account: fmt.Sprintf(
				"Assets:%s:%s",
				"ING",
				account.Name,
			),
			Metadata: []parser.MetadataField{{
				Key:   "tink_id",
				Value: fmt.Sprintf(`"%s"`, account.ID),
			}},
		})
	}

	for _, transaction := range fetchedTransactions {
		if _, ok := transactionsByTinkID[fmt.Sprintf(`"%s"`, transaction.ID)]; ok {
			continue
		}

		unscaled, err := strconv.Atoi(transaction.Amount.Value.UnscaledValue)
		if err != nil {
			return err
		}
		scale, err := strconv.Atoi(transaction.Amount.Value.Scale)
		if err != nil {
			return err
		}

		ledger = append(ledger, parser.Transaction{
			Date:      transaction.Dates.Booked,
			Type:      "!",
			Payee:     transaction.Descriptions.Display,
			Narration: transaction.Descriptions.Original,
			Metadata: []parser.MetadataField{{
				Key:   "tink_id",
				Value: fmt.Sprintf(`"%s"`, transaction.ID),
			}},
			Postings: []parser.Posting{{
				Account: "Assets:ING:Girokonto",
				Units: parser.Amount{
					Decimal: decimal.New(
						int64(unscaled),
						int32(scale*-1),
					),
					Currency: transaction.Amount.CurrencyCode,
				},
			}},
		})
	}

	return ledger.Save("./bla.beancount")
}

func fetchTransactions(
	accessToken string,
) ([]tink.Transaction, error) {
	client, err := tink.NewClient()
	if err != nil {
		return nil, err
	}

	transactions, err := client.Transactions(accessToken, "")
	if err != nil {
		return nil, fmt.Errorf("failed getting transactions: %w", err)
	}
	for transactions.NextPageToken != "" {
		res, err := client.Transactions(accessToken, transactions.NextPageToken)
		if err != nil {
			return nil, fmt.Errorf("failed getting transactions: %w", err)
		}
		transactions.Transactions = append(transactions.Transactions, res.Transactions...)
		transactions.NextPageToken = res.NextPageToken
	}

	return transactions.Transactions, nil
}

func fetchAccounts(
	accessToken string,
) ([]tink.Account, error) {
	client, err := tink.NewClient()
	if err != nil {
		return nil, err
	}

	accounts, err := client.Accounts(accessToken, "")
	if err != nil {
		return nil, fmt.Errorf("failed getting accounts: %w", err)
	}
	for accounts.NextPageToken != "" {
		res, err := client.Accounts(accessToken, accounts.NextPageToken)
		if err != nil {
			return nil, fmt.Errorf("failed getting accounts: %w", err)
		}
		accounts.Accounts = append(accounts.Accounts, res.Accounts...)
		accounts.NextPageToken = res.NextPageToken
	}

	return accounts.Accounts, nil
}
