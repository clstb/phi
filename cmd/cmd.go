package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Core(ctx *cli.Context) (pb.CoreClient, error) {
	configPath := ctx.String("config")
	config, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	coreHost := ctx.String("core-host")

	conn, err := grpc.Dial(
		coreHost,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				ctx = metadata.AppendToOutgoingContext(
					ctx,
					"authorization",
					fmt.Sprintf("Bearer %s", config.AccessToken),
				)
				return invoker(ctx, method, req, reply, cc, opts...)
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return pb.NewCoreClient(conn), nil
}

func Auth(ctx *cli.Context) (pb.AuthClient, error) {
	authHost := ctx.String("auth-host")

	conn, err := grpc.Dial(authHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewAuthClient(conn), nil
}

func renderTree(
	tree treeprint.Tree,
	accounts fin.Accounts,
	sum map[string]fin.Amounts,
) []byte {
	var amounts fin.Amounts
	for _, v := range sum {
		amounts = append(amounts, v...)
	}
	currencies := amounts.Sum().Currencies()

	s := ""
	for _, currency := range currencies {
		s += "\t" + currency
	}
	tree.SetValue(s)

	re := regexp.MustCompile("^(Income|Expenses|Equity)")
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
		invert := re.MatchString(account.Name)
		for _, currency := range currencies {
			amount := sum[account.ID.String()].ByCurrency(currency)
			s += "\t" + amount.ColorRaw(invert)
		}
		branch.SetValue(s)
	}

	return tree.Bytes()
}
