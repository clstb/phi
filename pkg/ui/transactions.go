package ui

import (
	"log"

	"github.com/clstb/phi/pkg/fin"
	"github.com/gdamore/tcell/v2"
	"github.com/gofrs/uuid"
	"github.com/rivo/tview"
)

func (u *UI) selectedTransaction() (fin.Transaction, int) {
	row, _ := u.tt.GetSelection()
	if row == 0 {
		return fin.Transaction{}, 0
	}
	return u.transactions[row-1], row - 1
}

func (u *UI) renderTransactions() {
	u.tt.Clear()

	header := []string{
		"Date",
		"Entity",
		"Reference",
	}

	for column, field := range header {
		u.tt.SetCell(0, column, &tview.TableCell{
			Color: tcell.ColorYellow,
			Text:  field,
		})
	}

	for row, transaction := range u.transactions {
		color := tcell.ColorWhite
		columns := []string{
			transaction.Date.Format("2006-01-02"),
			transaction.Entity,
			transaction.Reference,
		}

		if transaction.ID == uuid.Nil {
			color = tcell.ColorGrey
		}

		for column, field := range columns {
			u.tt.SetCell(row+1, column, &tview.TableCell{
				Text:  field,
				Color: color,
			})
		}
	}
}

func (u *UI) handlerTransactions() {
	u.tt.SetDoneFunc(func(key tcell.Key) {
		u.app.SetRoot(u.m, true)
	})
	u.tt.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}

		u.renderPostings()
		u.main.AddItem(u.side, 0, 1, false)
		u.app.SetFocus(u.pt)
	})
	u.tt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyRune {
			return event
		}

		transaction, row := u.selectedTransaction()
		if row == 0 {
			return event
		}

		switch event.Rune() {
		case 'd':
			if transaction.ID != uuid.Nil {
				_, err := u.core.DeleteTransaction(
					u.ctx,
					transaction.PB(),
				)
				if err != nil {
					log.Fatal(err)
				}
				u.transactions = append(
					u.transactions[:row],
					u.transactions[row+1:]...,
				)
				u.tt.Select(row, 0)
				u.render()
			}
			return nil
		case 's':
			if transaction.ID == uuid.Nil {
				transactionPB, err := u.core.CreateTransaction(
					u.ctx,
					transaction.PB(),
				)
				if err != nil {
					log.Fatal(err)
				}

				transaction, err := fin.TransactionFromPB(transactionPB)
				if err != nil {
					log.Fatal(err)
				}
				u.transactions[row] = transaction
				u.render()
			}
			return nil
		default:
			return event
		}
	})
}
