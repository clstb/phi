package create

import (
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
		Label:   "Entity",
		Default: "Username",
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

	transaction := fin.Transaction{}
	transaction.Date = date
	transaction.Entity = entity
	transaction.Reference = reference

	pSelect := promptui.Select{
		Label: "Select action",
		Items: []string{"Done", "Add posting"},
	}
	for {
		_, action, err := pSelect.Run()
		if err != nil {
			return err
		}
		if action == "Done" {
			break
		}
		posting, err := postingPrompt(accounts)
		if err != nil {
			return err
		}
		transaction.Postings = append(transaction.Postings, posting)
	}

	stream, err := core.CreateTransactions(ctx.Context)
	if err != nil {
		return err
	}
	if err := stream.Send(transaction.PB()); err != nil {
		return err
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func postingPrompt(accounts fin.Accounts) (fin.Posting, error) {
	accountNames := accounts.Names()

	getAccount := func() (fin.Account, error) {
		p := promptui.Select{
			Label:             "Select account",
			Items:             accountNames,
			StartInSearchMode: true,
			Searcher: func(s string, i int) bool {
				s = strings.ToLower(strings.TrimSpace(s))
				name := strings.ToLower(strings.TrimSpace(accountNames[i]))
				return fuzzy.Match(s, name)
			},
		}

		_, accountName, err := p.Run()
		if err != nil {
			return fin.Account{}, err
		}
		return accounts.ByName(accountName), nil
	}
	getAmount := func(label string) (fin.Amount, error) {
		p := promptui.Prompt{
			Label: "Units",
		}
		amountStr, err := p.Run()
		if err != nil {
			return fin.Amount{}, err
		}
		amount, err := fin.AmountFromString(amountStr)
		if err != nil {
			return fin.Amount{}, err
		}

		return amount, nil
	}

	account, err := getAccount()
	if err != nil {
		return fin.Posting{}, err
	}
	units, err := getAmount("Enter units")
	if err != nil {
		return fin.Posting{}, err
	}

	posting := fin.Posting{}
	posting.Account = account.ID
	posting.Units = units

	return posting, nil
}
