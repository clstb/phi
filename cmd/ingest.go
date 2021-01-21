package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/core/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/loader/csv"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Ingest(ctx *cli.Context) error {
	core, err := Core(ctx)
	if err != nil {
		return err
	}

	fp := ctx.Path("file")
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	config, err := config.Load(ctx.String("config"))
	if err != nil {
		return err
	}

	loader, err := csv.New(f, config, csv.WithSeperator(';'))
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

	hashes := make(map[string]struct{})
	for _, transaction := range transactionsPB.Data {
		hashes[transaction.Hash] = struct{}{}
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
	accountNames := accounts.Names()

	suggestions := []prompt.Suggest{
		{Text: "exit", Description: "Exit without saving"},
		{Text: "done", Description: "Exit and save"},
		{Text: "skip", Description: "Skip this transaction"},
		{Text: "help", Description: "Print help text"},
	}
	for _, s := range accountNames {
		suggestions = append(suggestions, prompt.Suggest{
			Text: s, Description: "Account",
		})
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

	transaction, amount, loadErr := loader.Load()
	fmt.Println(transaction, amount, loadErr)
	var toPush []fin.Transaction
L:
	for {
		if loadErr == io.EOF {
			break L
		}
		if loadErr != nil {
			return loadErr
		}

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
			transaction, amount, loadErr = loader.Load()
			continue
		}

		in := p.Input()
		in = strings.TrimSpace(in)

		switch {
		case reExit.MatchString(in):
			fmt.Println("Bye!")
			return nil
		case reDone.MatchString(in):
			break L
		case reSkip.MatchString(in):
			transaction, amount, loadErr = loader.Load()
		case reHelp.MatchString(in):
		case reAccAcc.MatchString(in):
			blocks := strings.Split(in, " ")
			from, ok := accounts.ByName(blocks[0])
			if !ok {
				fmt.Printf("Invalid account: %s\n", blocks[0])
				continue
			}

			to, ok := accounts.ByName(blocks[1])
			if !ok {
				fmt.Printf("Invalid account: %s\n", blocks[1])
				continue
			}

			postings := fin.Postings{
				fin.NewPosting(db.Posting{
					Account: from.ID,
					Units:   amount.Abs().Neg(),
				}),
				fin.NewPosting(db.Posting{
					Account: to.ID,
					Units:   amount.Abs(),
				}),
			}
			transaction.Postings = postings
			toPush = append(toPush, transaction)
			transaction, amount, loadErr = loader.Load()
		default:
			fmt.Println("Invalid command.")
		}
	}

	fmt.Println("Uploading transactions...")
	for _, transaction := range toPush {
		_, err = core.CreateTransaction(
			ctx.Context,
			transaction.PB(),
		)
		if err != nil {
			return err
		}
	}
	fmt.Println("Success!")

	return nil
}
