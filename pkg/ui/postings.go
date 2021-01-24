package ui

import (
	"github.com/clstb/phi/pkg/fin"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *UI) selectedPosting() (fin.Posting, int) {
	transaction, _ := u.selectedTransaction()
	row, _ := u.pt.GetSelection()
	return transaction.Postings[row-1], row - 1
}

func (u *UI) renderPostings() {
	u.pt.Clear()

	header := []string{
		"Account",
		"Units",
		"Cost",
		"Price",
	}
	for column, field := range header {
		u.pt.SetCell(0, column, &tview.TableCell{
			Color: tcell.ColorYellow,
			Text:  field,
		})
	}

	transaction, _ := u.selectedTransaction()
	for row, posting := range transaction.Postings {
		color := tcell.ColorWhite
		account, _ := u.accounts.ById(posting.Account.String())
		columns := []string{
			account.Name,
			posting.Units.StringRaw(),
			posting.Cost.StringRaw(),
			posting.Price.StringRaw(),
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
		u.main.RemoveItem(u.side)
		u.app.SetFocus(u.tt)
	})
	u.pt.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}

		u.renderPostings()

		u.pfEdit()
		u.side.AddItem(u.pf, 0, 1, false)
		u.app.SetFocus(u.pf)
	})
	u.pt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyRune {
			return event
		}
		if event.Rune() != 'i' {
			return event
		}

		u.pfAdd()
		u.side.AddItem(u.pf, 0, 1, false)
		u.app.SetFocus(u.pf)

		return nil
	})
}
