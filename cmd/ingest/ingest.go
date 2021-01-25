package ingest

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/clstb/phi/cmd"
	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/ui"
	"github.com/urfave/cli/v2"
)

func Ingest(ctx *cli.Context) error {
	core, err := cmd.Core(ctx)
	if err != nil {
		return err
	}

	accountsPB, err := core.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{},
	)
	if err != nil {
		return err
	}

	accounts, err := fin.AccountsFromPB(accountsPB)
	if err != nil {
		return err
	}

	transactionsPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{},
	)
	if err != nil {
		return err
	}

	existingTransactions, err := fin.TransactionsFromPB(transactionsPB)
	if err != nil {
		return err
	}

	hashes := make(map[string]fin.Transaction)
	for _, transaction := range existingTransactions {
		hashes[transaction.Hash] = transaction
	}

	fp := ctx.Path("file")
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	cp := ctx.Path("config")
	c, err := config.Load(cp)
	if err != nil {
		return err
	}

	fc, err := c.ForFile(filepath.Base(f.Name()))
	if err != nil {
		return err
	}
	parse := parser(fc)

	r := csv.NewReader(f)
	r.Comma = ';'

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	var transactions fin.Transactions
	for _, record := range records {
		transaction, err := parse(record)
		if err != nil {
			return err
		}

		existingTransaction, ok := hashes[transaction.Hash]
		if ok {
			transactions = append(transactions, existingTransaction)
		} else {
			transactions = append(transactions, transaction)
		}

	}

	ui := ui.New(
		ctx.Context,
		transactions,
		accounts,
		core,
	)
	return ui.Run()
}

func parser(
	c config.FileConfig,
) func([]string) (fin.Transaction, error) {
	return func(s []string) (fin.Transaction, error) {
		amount, err := fin.AmountFromString(
			fmt.Sprintf("%s %s",
				s[c.Amount],
				s[c.Currency],
			),
			fin.AmountEU,
		)
		if err != nil {
			return fin.Transaction{}, err
		}

		date, err := time.Parse(c.DateFormat, s[c.Date])
		if err != nil {
			return fin.Transaction{}, err
		}

		hash := sha256.New()
		_, err = hash.Write([]byte(strings.Join(s, "")))
		if err != nil {
			return fin.Transaction{}, err
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))

		posting := fin.Posting{}
		posting.Units = amount

		transaction := fin.Transaction{}
		transaction.Date = date
		transaction.Hash = hashStr
		transaction.Entity = s[c.Entity]
		transaction.Reference = s[c.Reference]
		transaction.Postings = append(
			transaction.Postings,
			posting,
		)
		return transaction, nil
	}
}
