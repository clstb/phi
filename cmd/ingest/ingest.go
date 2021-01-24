package ingest

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/clstb/phi/cmd"
	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

func Ingest(ctx *cli.Context) error {
	fp := ctx.Path("file")
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	cp := ctx.Path("config")
	c, err := config.Load(cp)
	if err != nil {
		return err
	}

	fc, err := c.ForFile(filepath.Base(f.Name()))
	if err != nil {
		return err
	}
	parse := parser(fc)

	r := csv.NewReader(f)
	r.Comma = ';'

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	var transactions fin.Transactions
	var amounts fin.Amounts
	for _, record := range records {
		transaction, amount, err := parse(record)
		if err != nil {
			return err
		}
		transactions = append(transactions, transaction)
		amounts = append(amounts, amount)
	}

	core, err := cmd.Core(ctx)
	if err != nil {
		return err
	}

	accountsPB, err := core.GetAccounts(
		ctx.Context,
		&pb.AccountsQuery{},
	)
	if err != nil {
		return err
	}

	accounts, err := fin.AccountsFromPB(accountsPB)
	if err != nil {
		return err
	}

	transactionsPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{},
	)
	if err != nil {
		return err
	}

	hashes := make(map[string]struct{})
	for _, transaction := range transactionsPB.Data {
		hashes[transaction.Hash] = struct{}{}
	}

	app := tview.NewApplication()

	main := tview.NewFlex()
	side := tview.NewFlex().SetDirection(tview.FlexRow)

	tt := tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	tt.SetBorder(true).SetTitle("Transactions")
	tt.SetEvaluateAllRows(true)
	currentTransaction := func() (fin.Transaction, int) {
		row, _ := tt.GetSelection()
		return transactions[row-1], row - 1
	}

	pt := tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	pt.SetBorder(true).SetTitle("Postings")
	pt.SetEvaluateAllRows(true)
	currentPosting := func() (fin.Posting, int) {
		transaction, _ := currentTransaction()
		row, _ := pt.GetSelection()
		return transaction.Postings[row-1], row - 1
	}

	m := tview.NewModal()
	m.SetText("Do you want to quit the application?")
	m.AddButtons([]string{"Quit & Save", "Quit", "Cancel"})
	m.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Quit":
			app.Stop()
		case "Cancel":
			app.SetRoot(main, true)
		case "Quit & Save":
			// TODO: the user should get feedback about transaction uploading and errors
			for _, transaction := range transactions {
				_, ok := hashes[transaction.Hash]
				if ok {
					continue
				}
				if !transaction.Balanced() {
					continue
				}
				if len(transaction.Postings) == 0 {
					continue
				}

				_, err := core.CreateTransaction(
					ctx.Context,
					transaction.PB(),
				)

				if err != nil {
					log.Fatal(err)
				}
			}
			app.Stop()
		}
	})

	inputAccount := tview.NewInputField()
	inputAccount.SetLabel("Account")
	inputAccount.SetAutocompleteFunc(func(currentText string) (entries []string) {
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

	inputUnits := tview.NewInputField()
	inputUnits.SetLabel("Units")

	inputCost := tview.NewInputField()
	inputCost.SetLabel("Cost")

	inputPrice := tview.NewInputField()
	inputPrice.SetLabel("Price")

	pf := tview.NewForm()
	pf.SetBorder(true)
	pf.AddFormItem(inputAccount)
	pf.AddFormItem(inputUnits)
	pf.AddFormItem(inputCost)
	pf.AddFormItem(inputPrice)

	parsePosting := func() (fin.Posting, error) {
		account, ok := accounts.ByName(inputAccount.GetText())
		if !ok {
			inputAccount.SetText("Invalid account")
			return fin.Posting{}, fmt.Errorf("Invalid account")
		}

		units, err := fin.AmountFromString(inputUnits.GetText())
		if err != nil {

			inputUnits.SetText("Invalid format")
			return fin.Posting{}, fmt.Errorf("Invalid format")
		}

		cost, err := fin.AmountFromString(inputCost.GetText())
		if err != nil {
			inputCost.SetText("Invalid format")
			return fin.Posting{}, fmt.Errorf("Invalid format")
		}

		price, err := fin.AmountFromString(inputPrice.GetText())
		if err != nil {
			inputPrice.SetText("Invalid format")
			return fin.Posting{}, fmt.Errorf("Invalid format")
		}

		posting := fin.Posting{}
		posting.Account = account.ID
		posting.Units = units
		posting.Cost = cost
		posting.Price = price

		return posting, nil
	}
	pfEdit := func(pRow int) {
		transaction, tRow := currentTransaction()
		posting, pRow := currentPosting()

		account, ok := accounts.ById(posting.Account.String())
		if !ok {
			inputAccount.SetText("Invalid account")
		} else {
			inputAccount.SetText(account.Name)
		}
		inputUnits.SetText(posting.Units.String())
		inputCost.SetText(posting.Cost.String())
		inputPrice.SetText(posting.Price.String())

		pf.Clear(true)
		pf.SetTitle("Edit Posting")
		pf.AddFormItem(inputAccount)
		pf.AddFormItem(inputUnits)
		pf.AddFormItem(inputCost)
		pf.AddFormItem(inputPrice)
		pf.AddButton("Save", func() {
			posting, err := parsePosting()
			if err != nil {
				return
			}

			transaction.Postings[pRow] = posting
			transactions[tRow] = transaction

			renderPostings(pt, transaction.Postings, accounts)
			renderTransactions(tt, transactions, amounts, hashes)
			side.RemoveItem(pf)
			app.SetFocus(pt)
		})
	}
	pfAdd := func() {
		transaction, tRow := currentTransaction()
		amount := amounts[tRow]
		if len(transaction.Postings) == 0 {
			inputUnits.SetText(amount.String())
		} else {
			sum := transaction.Postings.Sum()

			var amounts fin.Amounts
			for _, v := range sum {
				amounts = append(amounts, v...)

			}
			amounts = amounts.Sum()
			var amountsStr []string
			for _, v := range amounts {
				amountsStr = append(amountsStr, v.Neg().String())
			}

			inputUnits.SetText(strings.Join(amountsStr, ";"))
		}

		inputAccount.SetText("")
		inputCost.SetText("")
		inputPrice.SetText("")

		pf.Clear(true)
		pf.SetTitle("Add Posting")
		pf.AddFormItem(inputAccount)
		pf.AddFormItem(inputUnits)
		pf.AddFormItem(inputCost)
		pf.AddFormItem(inputPrice)
		pf.AddButton("Save", func() {
			posting, err := parsePosting()
			if err != nil {
				return
			}

			transaction.Postings = append(transaction.Postings, posting)
			transactions[tRow] = transaction
			renderPostings(pt, transaction.Postings, accounts)
			renderTransactions(tt, transactions, amounts, hashes)
			side.RemoveItem(pf)
			app.SetFocus(pt)
		})
	}

	tt.SetDoneFunc(func(key tcell.Key) {
		app.SetRoot(m, true)
	})
	tt.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}

		transaction, _ := currentTransaction()
		renderPostings(pt, transaction.Postings, accounts)
		main.AddItem(side, 0, 1, false)
		app.SetFocus(pt)
	})

	pt.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyESC {
			return
		}
		main.RemoveItem(side)
		app.SetFocus(tt)
	})
	pt.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}

		pfEdit(row)
		side.AddItem(pf, 0, 1, false)
		app.SetFocus(pf)
	})
	pt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyRune {
			return event
		}
		if event.Rune() != 'i' {
			return event
		}

		pfAdd()
		side.AddItem(pf, 0, 1, false)
		app.SetFocus(pf)
		return nil
	})

	pf.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyESC {
			return event
		}

		side.RemoveItem(pf)
		app.SetFocus(pt)
		return nil
	})

	renderTransactions(tt, transactions, amounts, hashes)

	main.AddItem(tt, 0, 3, true)
	side.AddItem(pt, 0, 1, true)

	return app.SetRoot(main, true).Run()
}

func renderTransactions(
	t *tview.Table,
	transactions fin.Transactions,
	amounts fin.Amounts,
	hashes map[string]struct{},
) {
	t.Clear()

	header := []string{
		"Date",
		"Amount",
		"Currency",
		"Entity",
		"Reference",
	}
	for column, field := range header {
		t.SetCell(0, column, &tview.TableCell{
			Color: tcell.ColorYellow,
			Text:  field,
		})
	}

	for row, transaction := range transactions {
		color := tcell.ColorWhite
		columns := []string{
			transaction.Date.Format("2006-01-02"),
			amounts[row].StringRaw(),
			amounts[row].Currency,
			transaction.Entity,
			transaction.Reference,
		}

		_, ok := hashes[transaction.Hash]
		if ok {
			color = tcell.ColorGray
		}
		if transaction.Balanced() && len(transaction.Postings) != 0 {
			color = tcell.ColorGreen
		}
		if !transaction.Balanced() {
			color = tcell.ColorRed
		}

		for column, field := range columns {
			t.SetCell(row+1, column, &tview.TableCell{
				Text:  field,
				Color: color,
			})
		}
	}

}

func renderPostings(
	t *tview.Table,
	postings fin.Postings,
	accounts fin.Accounts,
) {
	t.Clear()

	header := []string{
		"Account",
		"Units",
		"Cost",
		"Reference",
	}
	for column, field := range header {
		t.SetCell(0, column, &tview.TableCell{
			Color: tcell.ColorYellow,
			Text:  field,
		})
	}

	for row, posting := range postings {
		color := tcell.ColorWhite
		account, _ := accounts.ById(posting.Account.String())
		columns := []string{
			account.Name,
			posting.Units.StringRaw(),
			posting.Cost.StringRaw(),
			posting.Price.StringRaw(),
		}

		for column, field := range columns {
			t.SetCell(row+1, column, &tview.TableCell{
				Text:  field,
				Color: color,
			})
		}
	}
}

func parser(
	c config.FileConfig,
) func([]string) (fin.Transaction, fin.Amount, error) {
	return func(s []string) (fin.Transaction, fin.Amount, error) {
		amount, err := fin.AmountFromString(
			fmt.Sprintf("%s %s",
				s[c.Amount],
				s[c.Currency],
			),
			fin.AmountEU,
		)
		if err != nil {
			return fin.Transaction{}, fin.Amount{}, err
		}

		date, err := time.Parse(c.DateFormat, s[c.Date])
		if err != nil {
			return fin.Transaction{}, fin.Amount{}, err
		}

		hash := sha256.New()
		_, err = hash.Write([]byte(strings.Join(s, "")))
		if err != nil {
			return fin.Transaction{}, fin.Amount{}, err
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))

		transaction := fin.Transaction{}
		transaction.Date = date
		transaction.Hash = hashStr
		transaction.Entity = s[c.Entity]
		transaction.Reference = s[c.Reference]
		return transaction, amount, nil

	}
}
