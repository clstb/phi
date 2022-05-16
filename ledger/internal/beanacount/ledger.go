package beanacount

import (
	"bufio"
	"fmt"
	"github.com/clstb/phi/ledger/internal/config"
	"os"
)

type Ledger []fmt.Stringer

func (l *Ledger) Transactions() Transactions {
	var transactions Transactions
	for _, v := range *l {
		transaction, ok := v.(Transaction)
		if !ok {
			continue
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func (l Ledger) Opens() Opens {
	var opens Opens
	for _, v := range l {
		open, ok := v.(Open)
		if !ok {
			continue
		}
		opens = append(opens, open)
	}
	return opens
}

func GetFilePath(username string) string {
	return fmt.Sprintf("%s/%s.beancount", config.DataDirPath, username)
}
func (l *Ledger) PersistLedger(username string) error {
	// create file if not exists
	filePath := GetFilePath(username)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)
	defer datawriter.Flush()
	defer file.Close()

	// add default accounts
	datawriter.WriteString(config.DefaultOperatingCurrency)
	datawriter.WriteString(config.UnassignedIncomeAccount)
	datawriter.WriteString(config.UnassignedExpensesAccount)

	for _, i := range *l {
		_, _ = datawriter.WriteString(i.String() + "\n")

	}
	return nil
}
