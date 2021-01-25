package ui

import (
	"context"
	"strings"

	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
)

type UI struct {
	ctx context.Context

	transactions fin.Transactions
	accounts     fin.Accounts
	core         pb.CoreClient

	app  *tview.Application
	main *tview.Flex
	side *tview.Flex

	tt   *tview.Table
	tfid *tview.InputField
	tfie *tview.InputField
	tfir *tview.InputField
	tf   *tview.Form

	pt   *tview.Table
	pfia *tview.InputField
	pfiu *tview.InputField
	pfic *tview.InputField
	pfip *tview.InputField
	pf   *tview.Form

	m *tview.Modal
}

func New(
	ctx context.Context,
	transactions fin.Transactions,
	accounts fin.Accounts,
	core pb.CoreClient,
) *UI {
	app := tview.NewApplication()

	main := tview.NewFlex()
	side := tview.NewFlex().SetDirection(tview.FlexRow)

	tt := tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	tt.SetBorder(true).SetTitle("Transactions")
	tt.SetEvaluateAllRows(true)

	tfid := tview.NewInputField()
	tfid.SetLabel("Date")

	tfie := tview.NewInputField()
	tfie.SetLabel("Entity")

	tfir := tview.NewInputField()
	tfir.SetLabel("Reference")

	tf := tview.NewForm()
	tf.SetBorder(true)

	pt := tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	pt.SetBorder(true).SetTitle("Postings")
	pt.SetEvaluateAllRows(true)

	pfia := tview.NewInputField()
	pfia.SetLabel("Account")
	pfia.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}

		for _, account := range accounts {
			if fuzzy.Match(strings.ToLower(currentText), strings.ToLower(account.Name)) {
				entries = append(entries, account.Name)
			}
		}

		return
	})

	pfiu := tview.NewInputField()
	pfiu.SetLabel("Units")

	pfic := tview.NewInputField()
	pfic.SetLabel("Cost")

	pfip := tview.NewInputField()
	pfip.SetLabel("Price")

	pf := tview.NewForm()
	pf.SetBorder(true)

	m := tview.NewModal()
	m.SetText("Do you want to quit the application?")
	m.AddButtons([]string{"Quit", "Cancel"})

	side.AddItem(pt, 0, 1, true)
	main.AddItem(tt, 0, 3, true)
	main.AddItem(side, 0, 1, true)

	ui := &UI{
		ctx: ctx,

		transactions: transactions,
		accounts:     accounts,
		core:         core,

		app:  app,
		main: main,
		side: side,

		tt:   tt,
		tfid: tfid,
		tfie: tfie,
		tfir: tfir,
		tf:   tf,

		pt:   pt,
		pfia: pfia,
		pfiu: pfiu,
		pfic: pfic,
		pfip: pfip,
		pf:   pf,

		m: m,
	}

	ui.handlerTransactions()
	ui.handlerTransactionForm()
	ui.handlerPostings()
	ui.handlerPostingForm()
	ui.handlerModal()

	ui.renderTransactions()

	return ui
}

func (u *UI) Run() error {
	return u.app.SetRoot(u.main, true).Run()
}
