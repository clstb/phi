package cmd

import (
	"strings"

	"github.com/clstb/phi/pkg/db"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
	"google.golang.org/grpc"
)

func getClient(ctx *cli.Context) (pb.CoreClient, error) {
	apiHost := ctx.String("api-host")

	conn, err := grpc.Dial(apiHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewCoreClient(conn), nil
}

func renderTree(
	tree treeprint.Tree,
	accounts fin.Accounts,
	sum map[string]db.Amounts,
) []byte {
	var amounts db.Amounts
	for _, v := range sum {
		amounts = append(amounts, v...)
	}
	currencies := amounts.Sum().Currencies()

	s := ""
	for _, currency := range currencies {
		s += "\t" + currency
	}
	tree.SetValue(s)

	m := make(map[string]treeprint.Tree)
	for _, account := range accounts {
		path := strings.Split(account.Name, ":")
		branch := tree
		for _, s := range path {
			v, ok := m[s]
			if ok {
				branch = v
				continue
			}
			branch = branch.AddMetaBranch(s, "\t")
			m[s] = branch
		}
		s := ""
		for _, currency := range currencies {
			amount := sum[account.ID.String()].ByCurrency(currency)
			s += "\t" + amount.StringRaw()
		}
		branch.SetValue(s)
	}

	return tree.Bytes()
}
