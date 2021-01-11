package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Ingest(ctx *cli.Context) error {
	fp := ctx.Path("file")
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewCoreClient(conn)

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
	fmt.Println(accounts)

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

	var transactions []*fin.Transaction
	i := 0
	for i < len(records) {
		record := records[i]
		fmt.Println(record)

		in := p.Input()
		in = strings.TrimSpace(in)

		switch {
		case reExit.MatchString(in):
			fmt.Println("Bye!")
			return nil
		case reDone.MatchString(in):
			fmt.Println("Uploading transactions...")
			for _, transaction := range transactions {
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

			postings := fin.Postings{
				Data: []fin.Posting{
					{
						Account: from.Id,
						Units:   amount.Abs().Neg(),
					},
					{
						Account: to.Id,
						Units:   amount.Abs(),
					},
				},
			}

			transactions = append(transactions, &fin.Transaction{
				Date:      date.Format("2006-01-02"),
				Entity:    record[2],
				Reference: record[4],
				Postings:  postings,
			})
			i++
		default:
			fmt.Println("Invalid command.")
			continue
		}
	}
	return nil
}
