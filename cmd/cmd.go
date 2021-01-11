package cmd

import (
	"strings"

	"github.com/clstb/phi/pkg/fin"
	"github.com/xlab/treeprint"
)

func renderTree(
	tree treeprint.Tree,
	accounts fin.Accounts,
	sum fin.Sum,
	sumByCurrency fin.SumCurrency,
) []byte {
	s := ""
	for currency := range sumByCurrency {
		s += "\t" + currency

	}
	tree.SetValue(s)

	m := make(map[string]treeprint.Tree)
	for _, account := range accounts.Data {
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
		for currency := range sumByCurrency {
			amount, ok := sum[account.Id][currency]
			if !ok {
				s += "\t0"
			} else {
				s += "\t" + amount.StringRaw()
			}
		}
		branch.SetValue(s)
	}

	return tree.Bytes()
}

/*

func parseSums(
	accounts *pb.Accounts,
	sums *pb.Sums,
) map[string]string {
	m := make(map[string]string)
	r := regexp.MustCompilePOSIX("^(Income|Expenses|Equity)")

	for _, account := range accounts.Data {
		sum := sums.Data[account.Id]
		figure := parseFigure(sum, 2)

		m[account.Id] = colorFigure(
			figure,
			r.MatchString(account.Fullname),
		)
	}

	return m
}

func colorFigure(figure string, invert bool) string {
	red := false
	if strings.HasPrefix(figure, "-") {
		red = true
	}
	if invert {
		red = !red
	}
	if red {
		return color.RedString(figure)
	}

	return color.GreenString(figure)
}

func parseFigure(figure int64, percision int) string {
	figureRaw := strconv.Itoa(int(figure))
	var s string
	for i, r := range figureRaw {
		if i == len(figureRaw)-percision {
			s += "."
		}
		s += string(r)
	}
	return s + "€"
}
*/
