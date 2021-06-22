package cmd

import (
	"fmt"
	"strings"

	"github.com/clstb/phi/pkg/nordigen"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func Connect(ctx *cli.Context) error {
	token := ctx.String("nordigen-token")
	nordigen := nordigen.NewClient(token)

	banks, err := nordigen.GetBanks("de")
	if err != nil {
		return err
	}

	var bankNames []string
	for _, bank := range banks {
		bankNames = append(bankNames, bank.Name)
	}

	p := promptui.Select{
		Label:             "Select Bank",
		Items:             bankNames,
		StartInSearchMode: true,
		Searcher: func(s string, i int) bool {
			s = strings.TrimSpace(s)
			name := strings.TrimSpace(bankNames[i])
			return fuzzy.Match(s, name)
		},
	}

	i, result, err := p.Run()
	if err != nil {
		return err
	}

	eua, err := nordigen.CreateEndUserAgreement(
		"365",
		"test",
		banks[i].ID,
	)
	if err != nil {
		return err
	}
	fmt.Println(result, eua)

	return nil
}
