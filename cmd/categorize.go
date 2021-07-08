package cmd

import (
	"fmt"
	"strings"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
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

	classifier, err := bayesian.NewClassifierFromFile("./phi.classifier")
	if err != nil {
		return err
	}
	classes := classifier.Classes

	categorized := fin.Transactions{}
L:
	for _, transaction := range transactions {
		_, inx, _ := classifier.LogScores(strings.Split(
			transaction.Entity,
			" ",
		))
		class := strings.Split(string(classes[inx]), " ")
		fromID, toID := class[0], class[1]

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
				accountsById[fromID].Name,
				accountsById[toID].Name,
			),
		)

		var fromName string
		for {
			accountNames := append([]string{"Add account", "Done"}, accounts.Names()...)
			pSelect := promptui.Select{
				Label:             "From",
				Items:             accountNames,
				StartInSearchMode: true,
				Searcher: func(s string, i int) bool {
					s = strings.ToLower(strings.TrimSpace(s))
					name := strings.ToLower(strings.TrimSpace(accountNames[i]))
					return fuzzy.Match(s, name)
				},
			}
			_, fromName, err = pSelect.Run()
			if err != nil {
				return err
			}

			if fromName == "Add account" {
				p := promptui.Prompt{
					Label: "Account name",
				}
				accountName, err := p.Run()
				if err != nil {
					return err
				}

				accountPB, err := core.CreateAccount(
					ctx.Context,
					&pb.Account{
						Name: accountName,
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
			} else if fromName == "Done" {
				break L
			} else {
				break
			}
		}

		var toName string
		for {
			accountNames := append([]string{"Add account", "Done"}, accounts.Names()...)
			pSelect := promptui.Select{
				Label:             "To",
				Items:             accountNames,
				StartInSearchMode: true,
				Searcher: func(s string, i int) bool {
					s = strings.ToLower(strings.TrimSpace(s))
					name := strings.ToLower(strings.TrimSpace(accountNames[i]))
					return fuzzy.Match(s, name)
				},
			}
			_, toName, err = pSelect.Run()
			if err != nil {
				return err
			}

			if toName == "Add account" {
				p := promptui.Prompt{
					Label: "Account name",
				}
				accountName, err := p.Run()
				if err != nil {
					return err
				}

				accountPB, err := core.CreateAccount(
					ctx.Context,
					&pb.Account{
						Name: accountName,
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
			} else if toName == "Done" {
				break L
			} else {
				break
			}
		}

		transaction.From = accounts.ByName(fromName).ID
		transaction.To = accounts.ByName(toName).ID
		categorized = append(categorized, transaction)
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
