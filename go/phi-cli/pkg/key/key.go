package key

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Key struct {
	key.Binding
	Action func(tea.Model) (tea.Model, tea.Cmd)
}

func (k Key) Match(s string) bool {
	for _, key := range k.Keys() {
		if s == key {
			return true
		}
	}
	return false
}
