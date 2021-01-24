package ui

import "log"

func (u *UI) handlerModal() {
	u.m.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Quit":
			u.app.Stop()
		case "Cancel":
			u.app.SetRoot(u.main, true)
		case "Quit & Save":
			// TODO: the user should get feedback about transaction uploading and errors
			for _, transaction := range u.transactions {
				if !transaction.Balanced() {
					continue
				}
				if len(transaction.Postings) == 0 {
					continue
				}

				_, err := u.core.CreateTransaction(
					u.ctx,
					transaction.PB(),
				)

				if err != nil {
					log.Fatal(err)
				}
			}
			u.app.Stop()
		}
	})
}
