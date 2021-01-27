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

	app  *tview.Application // ui root
	main *tview.Flex        // main flex container
	side *tview.Flex        // side bar flex container

	tt   *tview.Table      // transactions table
	tfid *tview.InputField // transaction form input date
	tfie *tview.InputField // transaction form input entity
	tfir *tview.InputField // transaction form input reference
	tf   *tview.Form       // transaction form

	pt   *tview.Table      // postings table
	pfia *tview.InputField // posting form input account
	pfiu *tview.InputField // posting form input units
	pfic *tview.InputField // posting form input cost
	pfip *tview.InputField // posting form input price
	pf   *tview.Form       // posting form

	mq *tview.Modal // modal quit
	me *tview.Modal // modal error
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

	mq := tview.NewModal()
	mq.SetText("Do you want to quit the application?")
	mq.AddButtons([]string{"Quit", "Cancel"})

	me := tview.NewModal()
	me.AddButtons([]string{"Close"})

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

		mq: mq,
		me: me,
	}

	ui.handlerTransactions()
	ui.handlerTransactionForm()
	ui.handlerPostings()
	ui.handlerPostingForm()
	ui.handlerModalQuit()
	ui.handlerModalErr()

	ui.renderTransactions()

	return ui
}

func (u *UI) Run() error {
	return u.app.SetRoot(u.main, true).Run()
}
