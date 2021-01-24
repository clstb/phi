package ui

func (u *UI) handlerModal() {
	u.m.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Quit":
			u.app.Stop()
		case "Cancel":
			u.app.SetRoot(u.main, true)

		}
	})
}
