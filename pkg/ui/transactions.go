package ui

import (
	"github.com/clstb/phi/pkg/fin"
	"github.com/gdamore/tcell/v2"
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

		if !transaction.Balanced() {
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
}
