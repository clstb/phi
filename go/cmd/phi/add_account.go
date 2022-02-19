package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/pkg/ledger"
)

type addAccountModel struct {
	ledger ledger.Ledger
	input  textinput.Model
}

func newAddAccountModel() addAccountModel {
	input := textinput.New()
	input.Placeholder = "Account name"
	input.Focus()
	return addAccountModel{
		input: input,
	}
}

func (m addAccountModel) Init() tea.Cmd { return nil }

func (m addAccountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			v := m.input.Value()
			cmds = append(cmds, func() tea.Msg {
				return append(m.ledger, ledger.Open{
					Date:    "1970-01-01",
					Account: v,
				})
			})
			m.input.SetValue("")
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ledger.Ledger:
		m.ledger = msg
	}

	return m, tea.Batch(cmds...)
}

func (m addAccountModel) View() string {
	return m.input.View()
}
