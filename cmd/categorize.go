package cmd

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/navossoc/bayesian"
	"github.com/urfave/cli/v2"
)

func Categorize(ctx *cli.Context) error {
	core, err := Core(ctx)
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
	accountsById := accounts.ById()

	transactionsPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			AccountName: "^Uncategorized$",
		},
	)
	if err != nil {
		return err
	}
	transactions, err := fin.TransactionsFromPB(transactionsPB)
	if err != nil {
		return err
	}
	categorized := fin.Transactions{}

	classifier, err := bayesian.NewClassifierFromFile("./phi.classifier")
	if err != nil {
		return err
	}

	suggestions := []prompt.Suggest{
		{Text: "Done", Description: "Push transactions and exit"},
		{Text: "Add", Description: "Add new account"},
		{Text: "Skip", Description: "Skip this transaction"},
		{Text: "Exit", Description: "Exit with saving"},
	}
	for _, v := range accountsById {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        v.Name,
			Description: v.ID.String(),
		})
	}

	completer := func(in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == "" {
			return []prompt.Suggest{}
		}
		return prompt.FilterFuzzy(suggestions, w, true)
	}

	executor := func(in string) {
		in = strings.TrimSpace(in)
		fmt.Println(in)
	}

	p := prompt.New(
		executor,
		completer,
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.DarkGreen),
		prompt.OptionSelectedSuggestionBGColor(prompt.Black),
		prompt.OptionSelectedDescriptionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
	)

L:
	for i := 0; i < len(transactions); {
		transaction := transactions[i]

		_, inx, _ := classifier.LogScores(strings.Split(
			transaction.Entity,
			" ",
		))
		class := strings.Split(string(classifier.Classes[inx]), " ")

		from, to := accountsById[class[0]], accountsById[class[1]]
		transaction.From = from.ID
		transaction.To = to.ID

		if transaction.Debit {
			transaction.Units = transaction.Units.Neg()
			transaction.Cost = transaction.Cost.Neg()
			transaction.Price = transaction.Price.Neg()
		}

		fmt.Printf(
			"%s\n%s [%s] (%s) %s\n\n",
			transaction.Entity,
			transaction.Units.Color(false),
			transaction.Cost.Color(false),
			transaction.Price.Color(false),
			fmt.Sprintf(
				"%s => %s",
				from.Name,
				to.Name,
			),
		)

		in := p.Input()
		in = strings.TrimSpace(in)
		blocks := strings.Split(in, " ")

		switch {
		case blocks[0] == "Add":
			if len(blocks) != 2 {
				fmt.Printf("Please specify an account name!\n")
				continue
			}

			accountPB, err := core.CreateAccount(
				ctx.Context,
				&pb.Account{
					Name: blocks[1],
				},
			)
			if err != nil {
				return err
			}
			account, err := fin.AccountFromPB(accountPB)
			if err != nil {
				return err
			}

			accounts = append(accounts, account)
			suggestions = append(suggestions, prompt.Suggest{
				Text:        account.Name,
				Description: account.ID.String(),
			})
		case blocks[0] == "Done":
			break L
		case blocks[0] == "Exit":
			return nil
		case blocks[0] == "Skip":
			i++
		case blocks[0] == "":
			categorized = append(categorized, transaction)
			i++
		default:
			if len(blocks) != 2 {
				fmt.Printf("Please specify two account names!\n")
				break
			}

			from, to := accounts.ByName(blocks[0]), accounts.ByName(blocks[1])
			if from.Empty() {
				fmt.Printf("Account not found: %s\n", blocks[0])
				break
			}
			if to.Empty() {
				fmt.Printf("Account not found: %s\n", blocks[1])
				break
			}

			transaction.From, transaction.To = from.ID, to.ID
			categorized = append(categorized, transaction)
			i++
		}
	}

	_, err = core.UpdateTransactions(
		ctx.Context,
		categorized.PB(),
	)
	if err != nil {
		return err
	}

	return nil
}
