package cmd

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
)

func BalSheet(ctx *cli.Context) error {
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

	date := ctx.String("date")
	transactionsPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{
			To: date,
		},
	)

	transactions, err := fin.TransactionsFromPB(transactionsPB)
	if err != nil {
		return err
	}

	cleared, err := transactions.Clear(accounts)
	if err != nil {
		return err
	}

	sum, err := cleared.Sum()
	if err != nil {
		return err
	}

	amounts := fin.Amounts{}
	for _, v := range sum {
		amounts = append(amounts, v...)
	}
	amountsSum, err := amounts.Sum()
	if err != nil {
		return err
	}
	currencies := amountsSum.Currencies()

	tree := treeprint.New()
	tree.SetMetaValue("Balance Sheet")

	re := regexp.MustCompile("^(Equity|Liabilities|Assets)")
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	_, err = w.Write(renderTree(
		tree,
		accounts.FilterName(re),
		currencies,
		sum,
	))
	if err != nil {
		return err
	}

	re = regexp.MustCompile("^Equity")
	amounts = fin.Amounts{}
	for accountId, v := range sum {
		// these accounts always exist so we don't check for empty
		account := accounts.ById(accountId)
		if !re.MatchString(account.Name) {
			continue
		}
		amounts = append(amounts, v...)
	}

	nw, err := amounts.Sum()
	if err != nil {
		return err
	}

	var s string
	for _, amount := range nw {
		s += "\t" + amount.ColorRaw(false)
	}
	fmt.Fprintf(w, "\t\nNet Worth:%s\n", s)

	return w.Flush()
}
