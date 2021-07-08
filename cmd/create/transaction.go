package create

import (
	"database/sql"
	"strings"
	"time"

	"github.com/clstb/phi/cmd"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func Transaction(ctx *cli.Context) error {
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

	p := promptui.Prompt{
		Label: "Date",
		Validate: func(s string) error {
			_, err := time.Parse("2006-01-02", s)
			return err
		},
		Default: time.Now().Format("2006-01-02"),
	}
	dateStr, err := p.Run()
	if err != nil {
		return err
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	p = promptui.Prompt{
		Label: "Entity",
	}
	entity, err := p.Run()
	if err != nil {
		return err
	}

	p = promptui.Prompt{
		Label: "Reference",
	}
	reference, err := p.Run()
	if err != nil {
		return err
	}

	accountNames := accounts.Names()
	pSelect := promptui.Select{
		Label:             "Select account",
		Items:             accountNames,
		StartInSearchMode: true,
		Searcher: func(s string, i int) bool {
			s = strings.ToLower(strings.TrimSpace(s))
			name := strings.ToLower(strings.TrimSpace(accountNames[i]))
			return fuzzy.Match(s, name)
		},
	}

	_, fromName, err := pSelect.Run()
	if err != nil {
		return err
	}

	_, toName, err := pSelect.Run()
	if err != nil {
		return err
	}

	p = promptui.Prompt{
		Label: "Amount",
	}
	amountStr, err := p.Run()
	if err != nil {
		return err
	}
	amount, err := fin.AmountFromString(amountStr)
	if err != nil {
		return err
	}

	transaction := fin.Transaction{}
	transaction.Date = date
	transaction.Entity = entity
	if reference != "" {
		transaction.Reference = sql.NullString{
			String: reference,
			Valid:  true,
		}
	}
	transaction.From = accounts.ByName(fromName).ID
	transaction.To = accounts.ByName(toName).ID
	transaction.Units = amount.Abs()

	_, err = core.CreateTransactions(
		ctx.Context,
		&pb.Transactions{
			Data: []*pb.Transaction{transaction.PB()},
		},
	)
	if err != nil {
		return err
	}

	return nil
}
