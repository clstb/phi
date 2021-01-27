package ui

func (u *UI) mePrep(err error) {
	u.me.SetText(err.Error())
}

func (u *UI) handlerModalErr() {
	u.me.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Close":
			u.app.SetRoot(u.main, true)
		}
	})
}
