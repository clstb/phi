package ui

import (
	"fmt"
	"time"

	"github.com/clstb/phi/pkg/fin"
	"github.com/gdamore/tcell/v2"
)

func (u *UI) transactionFromForm() (fin.Transaction, error) {
	date, err := time.Parse("2006-01-02", u.tfid.GetText())
	if err != nil {
		u.tfid.SetText("Invalid date")
		return fin.Transaction{}, fmt.Errorf("Invalid date")
	}

	transaction := fin.Transaction{}
	transaction.Date = date
	transaction.Entity = u.tfie.GetText()
	transaction.Reference = u.tfir.GetText()

	return transaction, nil
}

func (u *UI) tfPrep(row int) {
	var transaction fin.Transaction
	if row > 0 {
		u.tf.SetTitle("Edit Transaction")
		transaction = u.transactions[row-1]
		u.tfid.SetText(transaction.Date.Format("2006-01-02"))
	} else {
		u.tf.SetTitle("Add Transaction")
		u.tfid.SetText("")
	}

	u.tfie.SetText(transaction.Entity)
	u.tfir.SetText(transaction.Reference)

	u.tf.Clear(true)
	u.tf.AddFormItem(u.tfid)
	u.tf.AddFormItem(u.tfie)
	u.tf.AddFormItem(u.tfir)
	u.tf.AddButton("Save", func() {
		transaction, err := u.transactionFromForm()
		if err != nil {
			return
		}

		if row > 0 {
			u.transactions = append(u.transactions, transaction)
		} else {
			u.transactions[row-1] = transaction
		}

		u.renderTransactions()
		u.side.RemoveItem(u.tf)
		u.app.SetFocus(u.tt)
	})
}

func (u *UI) handlerTransactionForm() {
	u.tf.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyESC {
			return event

		}

		u.side.RemoveItem(u.tf)
		u.app.SetFocus(u.tt)
		return nil
	})
}
