package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/internal/phi/beancount"
)

func LoadLedger(path string) tea.Cmd {
	return func() tea.Msg {
		var l beancount.Ledger
		filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			f, err := os.OpenFile(path, os.O_RDONLY, info.Mode())
			if err != nil {
				return err
			}

			l = append(l, beancount.Parse(f)...)
			return f.Close()
		})

		return l
	}
}

func SaveLedger(
	path string,
	ledger beancount.Ledger,
) tea.Cmd {
	return func() tea.Msg {
		transactionsByMonth := ledger.Transactions().ByMonth()
		for month, transactions := range transactionsByMonth {
			f, err := os.OpenFile(
				fmt.Sprintf("%s/transactions/%s.bean", path, month),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				os.ModePerm,
			)
			if err != nil {
				return err
			}
			defer f.Close()

			sort.Slice(transactions, func(i, j int) bool {
				return transactions[i].Date.Before(transactions[j].Date)
			})

			for _, transaction := range transactions {
				fmt.Fprint(f, transaction.String())
			}
		}

		f, err := os.OpenFile(
			fmt.Sprintf("%s/accounts.bean", path),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			os.ModePerm,
		)
		if err != nil {
			return err
		}
		defer f.Close()

		opens := ledger.Opens()
		sort.Slice(opens, func(i, j int) bool {
			return opens[i].Account < opens[j].Account
		})
		for _, open := range opens {
			fmt.Fprint(f, open.String())
		}

		return nil
	}
}
