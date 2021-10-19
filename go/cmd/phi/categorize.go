package main

import (
	"fmt"

	"github.com/clstb/phi/go/pkg/parser"
	"github.com/jbrukh/bayesian"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type classes map[bayesian.Class][]string

func (c classes) Classifier() *bayesian.Classifier {
	var classes []bayesian.Class
	for class := range c {
		classes = append(classes, class)
	}

	classifier := bayesian.NewClassifierTfIdf(classes...)
	for k, v := range c {
		classifier.Learn(v, k)
	}
	classifier.ConvertTermsFreqToTfIdf()
	return classifier
}

func Categorize(ctx *cli.Context) error {
	ledger, err := parser.Load("./bla.beancount")
	if err != nil {
		return err
	}

	opens := ledger.Opens()
	if len(opens) < 2 {
		return fmt.Errorf("ledger needs atleast 2 open accounts")
	}

	pClasses := classes{}
	nClasses := classes{}
	opensByAccount := map[string]parser.Open{}
	for _, open := range ledger.Opens() {
		opensByAccount[open.Account] = open
		pClasses[bayesian.Class(open.Account)] = []string{}
		nClasses[bayesian.Class(open.Account)] = []string{}
	}

	transactionsByPayee := map[string][]parser.Transaction{}
	for _, transaction := range ledger.Transactions() {
		if transaction.Balanced() {
			for _, posting := range transaction.Postings {
				if posting.Weight().IsNegative() {
					nClasses[bayesian.Class(posting.Account)] = append(
						nClasses[bayesian.Class(posting.Account)],
						transaction.NormalizedPayee()...,
					)
				} else {
					pClasses[bayesian.Class(posting.Account)] = append(
						pClasses[bayesian.Class(posting.Account)],
						transaction.NormalizedPayee()...,
					)
				}
			}

			continue
		}

		v, ok := transactionsByPayee[transaction.Payee]
		if !ok {
			transactionsByPayee[transaction.Payee] = []parser.Transaction{transaction}
		} else {
			transactionsByPayee[transaction.Payee] = append(v, transaction)
		}
	}

	pClassifier := pClasses.Classifier()
	nClassifier := nClasses.Classifier()

	newOpens := map[string]parser.Open{}
	i := 1
L:
	for _, v := range transactionsByPayee {
		fmt.Printf(
			"Transaction group %d of %d\n",
			i,
			len(transactionsByPayee),
		)

		var (
			fastForward bool
			account     string
		)

		for _, transaction := range v {
			fmt.Println(transaction.String())

			if !fastForward {
				var classifier *bayesian.Classifier
				if transaction.Negative() {
					classifier = pClassifier
				} else {
					classifier = nClassifier
				}

				_, i, _ := classifier.LogScores(
					transaction.NormalizedPayee(),
				)

				prompt := promptui.Prompt{
					Label:   "Balance with account",
					Default: string(classifier.Classes[i]),
				}

				account, err = prompt.Run()
				if err != nil {
					return err
				}

				if account == "done" {
					break L
				}
			}

			if !fastForward && len(v) > 1 {
				prompt := promptui.Prompt{
					Label: fmt.Sprintf(
						"Found %d transactions with same payee. Fast forward categorize for all?",
						len(v)-1,
					),
					IsConfirm: true,
				}

				_, err := prompt.Run()
				if err == nil {
					fastForward = true
				}
			}

			_, ok := opensByAccount[account]
			if !ok {
				newOpens[account] = parser.Open{
					Account: account,
					Date:    "1970-01-01",
				}
			}

			if transaction.Negative() {
				pClasses[bayesian.Class(account)] = append(
					pClasses[bayesian.Class(account)],
					transaction.NormalizedPayee()...,
				)
				pClassifier = pClasses.Classifier()
			} else {
				nClasses[bayesian.Class(account)] = append(
					nClasses[bayesian.Class(account)],
					transaction.NormalizedPayee()...,
				)
				nClassifier = nClasses.Classifier()

			}

			transaction.Postings = append(
				transaction.Postings,
				parser.Posting{
					Account: account,
				},
			)
			transaction.Type = "*"
			ledger[transaction.Index] = transaction
		}
		i++
	}

	fmt.Printf("Following new accounts need to open:\n\n")
	for _, open := range newOpens {
		fmt.Print(open.String())
	}

	return ledger.Save("./bla.beancount")
}
