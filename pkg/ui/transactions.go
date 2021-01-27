package ui

import (
	"github.com/clstb/phi/pkg/fin"
	"github.com/gdamore/tcell/v2"
	"github.com/gofrs/uuid"
	"github.com/rivo/tview"
)

func (u *UI) renderTransactions() {
	u.tt.Clear()

	header := []string{
		"Date",
		"Entity",
		"Reference",
	}

	for column, field := range header {
		u.tt.SetCell(0, column, &tview.TableCell{
			Color:         tcell.ColorYellow,
			Text:          field,
			NotSelectable: true,
		})
	}

	for row, transaction := range u.transactions {
		color := tcell.ColorWhite

		columns := []string{
			transaction.Date.Format("2006-01-02"),
			transaction.Entity,
			transaction.Reference,
		}

		if transaction.ID == uuid.Nil && !transaction.Balanced() {
			color = tcell.ColorLightSalmon
		}
		if transaction.ID == uuid.Nil && transaction.Balanced() {
			color = tcell.ColorLightGreen
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
	u.tt.SetSelectionChangedFunc(func(row, column int) {
		if row > len(u.transactions) || row < 1 {
			return
		}

		transaction := u.transactions[row-1]
		u.renderPostings(transaction)
	})
	u.tt.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyESC {
			return
		}
		u.app.SetRoot(u.mq, true)
	})
	u.tt.SetSelectedFunc(func(row, column int) {
		if row > len(u.transactions) {
			return
		}

		u.tfPrep(row)
		u.side.AddItem(u.tf, 0, 1, false)
		u.app.SetFocus(u.tf)
	})
	u.tt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.app.SetFocus(u.pt)
		case tcell.KeyRune:
		default:
			return event
		}

		switch event.Rune() {
		case 'i':
			u.tfPrep(-1)
			u.side.AddItem(u.tf, 0, 1, false)
			u.app.SetFocus(u.tf)
			return nil
		case 'd':
			tRow, _ := u.tt.GetSelection()
			if tRow > len(u.transactions) {
				return event
			}

			transaction := u.transactions[tRow-1]

			if transaction.ID != uuid.Nil {
				_, err := u.core.DeleteTransaction(
					u.ctx,
					transaction.PB(),
				)
				if err != nil {
					u.mePrep(err)
					u.app.SetRoot(u.me, true)
				}
				u.transactions = append(
					u.transactions[:tRow-1],
					u.transactions[tRow:]...,
				)
				u.tt.Select(tRow-1, 0)
				u.renderTransactions()
			}
			return nil
		case 's':
			tRow, _ := u.tt.GetSelection()
			if tRow > len(u.transactions) {
				return event
			}

			transaction := u.transactions[tRow-1]

			if transaction.ID == uuid.Nil {
				transactionPB, err := u.core.CreateTransaction(
					u.ctx,
					transaction.PB(),
				)
				if err != nil {
					u.mePrep(err)
					u.app.SetRoot(u.me, true)
				}

				transaction, err := fin.TransactionFromPB(transactionPB)
				if err != nil {
					u.mePrep(err)
					u.app.SetRoot(u.me, true)
				}

				u.transactions[tRow-1] = transaction
				u.renderTransactions()
			}
			return nil
		default:
			return event
		}
	})
}
