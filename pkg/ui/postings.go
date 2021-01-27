package ui

import (
	"github.com/clstb/phi/pkg/fin"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *UI) renderPostings(transaction fin.Transaction) {
	u.pt.Clear()

	header := []string{
		"Account",
		"Units",
		"Cost",
		"Price",
	}
	for column, field := range header {
		u.pt.SetCell(0, column, &tview.TableCell{
			Color:         tcell.ColorYellow,
			Text:          field,
			NotSelectable: true,
		})
	}

	for row, posting := range transaction.Postings {
		color := tcell.ColorWhite
		account, _ := u.accounts.ById(posting.Account.String())
		columns := []string{
			account.Name,
			posting.Units.String(),
			posting.Cost.String(),
			posting.Price.String(),
		}

		for column, field := range columns {
			u.pt.SetCell(row+1, column, &tview.TableCell{
				Text:  field,
				Color: color,
			})
		}
	}
}

func (u *UI) handlerPostings() {
	u.pt.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyESC {
			return
		}
		u.app.SetRoot(u.mq, true)
	})
	u.pt.SetSelectedFunc(func(row, column int) {
		tRow, _ := u.tt.GetSelection()
		transaction := u.transactions[tRow-1]

		if row > len(transaction.Postings) {
			return
		}

		u.pfPrep(row)
		u.side.AddItem(u.pf, 0, 1, false)
		u.app.SetFocus(u.pf)
	})
	u.pt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.app.SetFocus(u.tt)
		case tcell.KeyRune:
		default:
			return event
		}

		switch event.Rune() {
		case 'd':
			tRow, _ := u.tt.GetSelection()
			transaction := u.transactions[tRow-1]

			pRow, _ := u.pt.GetSelection()
			if pRow > len(transaction.Postings) {
				return event
			}

			transaction.Postings = append(
				transaction.Postings[:pRow-1],
				transaction.Postings[pRow:]...,
			)
			u.transactions[tRow-1] = transaction

			u.pt.Select(pRow-1, 0)

			u.renderPostings(transaction)
			u.renderTransactions()
		case 'i':
			u.pfPrep(-1)
			u.side.AddItem(u.pf, 0, 1, false)
			u.app.SetFocus(u.pf)
		default:
			return event
		}

		return nil
	})
}
