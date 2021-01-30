package csv

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/clstb/phi/cmd"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
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

	fp := ctx.Path("file")
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	var transactions fin.Transactions
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	from, to := records[0][0], records[0][1]
	if from == "" || to == "" {
		return fmt.Errorf("invalid file format: first entry needs to specify 2 accounts")
	}

	for _, record := range records {
		if record[0] != "" {
			from = record[0]
		}
		if record[1] != "" {
			to = record[1]
		}

		amount, err := fin.AmountFromString(
			record[5],
			fin.AmountEU,
		)
		if err != nil {
			return err
		}

		date, err := time.Parse(
			"2.1.2006",
			record[2],
		)
		if err != nil {
			return err
		}

		hash := sha256.New()
		_, err = hash.Write([]byte(strings.Join(record[2:], "")))
		if err != nil {
			return err
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))

		transaction := fin.Transaction{}
		transaction.Date = date
		transaction.Entity = record[3]
		transaction.Reference = record[4]
		transaction.Hash = hashStr

		accountFrom := accounts.ByName(from)
		if accountFrom.Empty() {
			fmt.Printf("Account not found %s: Creating...\n", from)
			accountPB, err := core.CreateAccount(
				ctx.Context,
				&pb.Account{
					Name: from,
				},
			)
			if err != nil {
				return err
			}

			accountFrom, err = fin.AccountFromPB(accountPB)
			if err != nil {
				return err
			}

			accounts = append(accounts, accountFrom)
		}

		accountTo := accounts.ByName(to)
		if accountTo.Empty() {
			fmt.Printf("Account not found %s: Creating...\n", to)
			accountPB, err := core.CreateAccount(
				ctx.Context,
				&pb.Account{
					Name: to,
				},
			)
			if err != nil {
				return err
			}

			accountTo, err = fin.AccountFromPB(accountPB)
			if err != nil {
				return err
			}

			accounts = append(accounts, accountTo)
		}

		posting := fin.Posting{}
		posting.Account = accountFrom.ID
		posting.Units = amount.Abs().Neg()
		transaction.Postings = append(transaction.Postings, posting)

		posting = fin.Posting{}
		posting.Account = accountTo.ID
		posting.Units = amount.Abs()
		transaction.Postings = append(transaction.Postings, posting)

		transactions = append(transactions, transaction)
	}

	stream, err := core.CreateTransactions(ctx.Context)
	if err != nil {
		return err
	}

	fmt.Printf("Uploading transactions...\n")
	for _, transaction := range transactions {
		err := stream.Send(transaction.PB())
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("Success!\n")

	return nil
}
