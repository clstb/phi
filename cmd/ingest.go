package cmd

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Ingest(ctx *cli.Context) error {
	client, err := getClient(ctx)
	if err != nil {
		return err
	}

	fp := ctx.Path("file")
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)
	r.Comma = ';'

	var transactions []fin.Transaction
	var amounts []fin.Amount
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		amount, err := fin.AmountFromString(
			record[7]+" "+record[8],
			fin.AmountEU,
		)
		if err != nil {
			return err
		}

		date, err := time.Parse("02.01.2006", record[1])
		if err != nil {
			return err
		}

		hash := sha256.New()
		_, err = hash.Write([]byte(strings.Join(record, "")))
		if err != nil {
			return err
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))

		transactions = append(
			transactions,
			fin.Transaction{
				Date:      date.Format("2006-01-02"),
				Entity:    record[2],
				Reference: record[4],
				Hash:      hashStr,
			},
		)
		amounts = append(amounts, amount)
	}

	transactionsPB, err := client.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			Fields: &pb.TransactionsQuery_Fields{
				Hash: true,
			},
			From: "-infinity",
			To:   "+infinity",
		},
	)
	if err != nil {
		return err
	}

	hashes := make(map[string]struct{})
	for _, transaction := range transactionsPB.Data {
		hashes[transaction.Hash] = struct{}{}
	}

	accounts, err := client.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{
			Fields: &pb.AccountsQuery_Fields{
				Name: true,
			},
		},
	)
	if err != nil {
		return err
	}

	suggestions := []prompt.Suggest{
		{Text: "exit", Description: "Exit without saving"},
		{Text: "done", Description: "Exit and save"},
		{Text: "skip", Description: "Skip this transaction"},
		{Text: "help", Description: "Print help text"},
	}
	var accountNames []string
	for _, account := range accounts.Data {
		suggestions = append(
			suggestions,
			prompt.Suggest{
				Text:        account.Name,
				Description: "Account",
			},
		)
		accountNames = append(accountNames, account.Name)
	}

	accountsRe := "(" + strings.Join(accountNames, "|") + ")"

	reExit := regexp.MustCompile("^exit$")
	reDone := regexp.MustCompile("^done$")
	reSkip := regexp.MustCompile("^skip$")
	reHelp := regexp.MustCompile("^help$")
	reAccAcc := regexp.MustCompile(fmt.Sprintf("^%s %s$", accountsRe, accountsRe))

	completer := func(in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == "" {
			return []prompt.Suggest{}
		}
		return prompt.FilterFuzzy(suggestions, w, true)
	}
	p := prompt.New(nil, completer)

	skipDuplicates := ctx.Bool("skip-duplicates")
	var toPush []fin.Transaction
	for i := 0; i < len(transactions); {
		transaction := transactions[i]
		amount := amounts[i]

		fmt.Printf(
			"Date:\t%s\nEntity:\t%s\nReference:\t%s\nAmount:\t%s\n\n",
			transaction.Date,
			transaction.Entity,
			transaction.Reference,
			amount.String(),
		)

		_, ok := hashes[transaction.Hash]
		if ok && skipDuplicates {
			fmt.Println("Found duplicate hash. Skipping...")
			i++
			continue
		}

		in := p.Input()
		in = strings.TrimSpace(in)

		switch {
		case reExit.MatchString(in):
			fmt.Println("Bye!")
			return nil
		case reDone.MatchString(in):
			fmt.Println("Uploading transactions...")
			for _, transaction := range toPush {
				pb, err := transaction.PB()
				if err != nil {
					return err
				}

				_, err = client.CreateTransaction(
					ctx.Context,
					pb,
				)
				if err != nil {
					return err
				}
			}
			fmt.Println("Success!")
			return nil
		case reSkip.MatchString(in):
			i++
		case reHelp.MatchString(in):
		case reAccAcc.MatchString(in):
			blocks := strings.Split(in, " ")
			from := accounts.Data[accounts.ByName[blocks[0]]]
			to := accounts.Data[accounts.ByName[blocks[1]]]

			postings := fin.Postings{
				Data: []fin.Posting{
					{
						Account: from.Id,
						Units:   amounts[i].Abs().Neg(),
					},
					{
						Account: to.Id,
						Units:   amounts[i].Abs(),
					},
				},
			}
			transaction.Postings = postings
			toPush = append(toPush, transaction)
			i++
		default:
			fmt.Println("Invalid command.")
			continue
		}
	}
	return nil
}
