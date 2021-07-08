package cmd

import (
	"strings"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/navossoc/bayesian"
	"github.com/urfave/cli/v2"
)

func Train(ctx *cli.Context) error {
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

	transactions, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			AccountName: "^(Income|Expenses|Assets|Equity|Liabilities)",
		},
	)
	if err != nil {
		return err
	}

	type ClassWithData struct {
		Class bayesian.Class
		Data  []string
	}

	m := map[bayesian.Class][]string{}
	for _, from := range accounts {
		for _, to := range accounts {
			m[bayesian.Class(from.ID.String()+" "+to.ID.String())] = []string{}
		}
	}

	for _, transaction := range transactions.Data {
		m[bayesian.Class(transaction.From+" "+transaction.To)] = append(
			m[bayesian.Class(transaction.From+" "+transaction.To)],
			strings.Split(transaction.Entity, " ")...,
		)
	}

	var classes []bayesian.Class
	for k := range m {
		classes = append(classes, k)
	}

	classifier := bayesian.NewClassifierTfIdf(classes...)
	for k, v := range m {
		classifier.Learn(v, k)
	}

	classifier.ConvertTermsFreqToTfIdf()

	if err := classifier.WriteToFile("./phi.classifier"); err != nil {
		return err
	}

	return nil
}
