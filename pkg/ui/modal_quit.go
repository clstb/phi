package ui

func (u *UI) handlerModalQuit() {
	u.mq.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Quit":
			u.app.Stop()
		case "Cancel":
			u.app.SetRoot(u.main, true)

		}
	})
}
