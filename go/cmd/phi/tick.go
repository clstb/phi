package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tick time.Time

func ticker() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tick(t)
	})
}
