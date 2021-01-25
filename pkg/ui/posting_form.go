package ui

import (
	"fmt"

	"github.com/clstb/phi/pkg/fin"
	"github.com/gdamore/tcell/v2"
)

func (u *UI) postingFromForm() (fin.Posting, error) {
	account, ok := u.accounts.ByName(u.pfia.GetText())
	if !ok {
		u.pfia.SetText("Invalid account")
		return fin.Posting{}, fmt.Errorf("Invalid account")
	}

	units, err := fin.AmountFromString(u.pfiu.GetText())
	if err != nil {

		u.pfiu.SetText("Invalid format")
		return fin.Posting{}, fmt.Errorf("Invalid format")
	}

	cost, err := fin.AmountFromString(u.pfic.GetText())
	if err != nil {
		u.pfic.SetText("Invalid format")
		return fin.Posting{}, fmt.Errorf("Invalid format")
	}

	price, err := fin.AmountFromString(u.pfip.GetText())
	if err != nil {
		u.pfip.SetText("Invalid format")
		return fin.Posting{}, fmt.Errorf("Invalid format")
	}

	posting := fin.Posting{}
	posting.Account = account.ID
	posting.Units = units
	posting.Cost = cost
	posting.Price = price

	return posting, nil

}

func (u *UI) pfPrep(row int) {
	tRow, _ := u.tt.GetSelection()
	transaction := u.transactions[tRow-1]

	var posting fin.Posting
	if row > 0 {
		u.pf.SetTitle("Edit Posting")
		posting = transaction.Postings[row-1]
	} else {
		u.pf.SetTitle("Add Posting")
		if len(transaction.Postings) == 1 {
			posting.Units = transaction.Postings[0].Units.Neg()
		}
	}

	account, ok := u.accounts.ById(posting.Account.String())
	if !ok {
		u.pfia.SetText("")
	} else {
		u.pfia.SetText(account.Name)
	}
	u.pfiu.SetText(posting.Units.String())
	u.pfic.SetText(posting.Cost.String())
	u.pfip.SetText(posting.Price.String())

	u.pf.Clear(true)
	u.pf.AddFormItem(u.pfia)
	u.pf.AddFormItem(u.pfiu)
	u.pf.AddFormItem(u.pfic)
	u.pf.AddFormItem(u.pfip)
	u.pf.AddButton("Save", func() {
		posting, err := u.postingFromForm()
		if err != nil {
			return
		}

		if row > 0 {
			transaction.Postings = append(transaction.Postings, posting)
		} else {
			transaction.Postings[row-1] = posting
		}

		u.transactions[tRow-1] = transaction

		u.renderTransactions()
		u.renderPostings(transaction)

		u.side.RemoveItem(u.pf)
		u.app.SetFocus(u.pt)
	})
}

func (u *UI) handlerPostingForm() {
	u.pf.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyESC {
			return event
		}

		u.side.RemoveItem(u.pf)
		u.app.SetFocus(u.pt)
		return nil
	})
}
