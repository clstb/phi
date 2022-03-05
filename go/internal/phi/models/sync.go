package models

import (
	"github.com/charmbracelet/bubbles/help"
	k "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/internal/phi/beancount"
	"github.com/clstb/phi/go/internal/phi/commands"
	"github.com/clstb/phi/go/internal/phi/key"
	"github.com/clstb/phi/go/internal/phi/state"
	"github.com/clstb/phi/go/pkg/client"
)

type Sync struct {
	ledger beancount.Ledger
	client *client.Client
	help   help.Model
}

func NewSync(client *client.Client) Sync {
	return Sync{
		client: client,
		help:   help.NewModel(),
	}
}

func (m Sync) Init() tea.Cmd {
	return commands.Sync(m.ledger, m.client)
}

func (m Sync) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		for _, key := range authKeys.slice() {
			if key.Match(msg.String()) {
				newModel, cmd := key.Action(m)
				m, cmds = newModel.(Sync), append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Sync) View() string {
	return m.help.View(syncKeys)
}

type syncKeymap struct {
	Help   key.Key
	Quit   key.Key
	Return key.Key
}

func (km syncKeymap) ShortHelp() []k.Binding {
	return []k.Binding{km.Help.Binding, km.Quit.Binding, km.Return.Binding}
}

func (km syncKeymap) FullHelp() [][]k.Binding {
	return [][]k.Binding{
		{km.Help.Binding, km.Quit.Binding, km.Return.Binding},
	}
}

func (km syncKeymap) slice() []key.Key {
	return []key.Key{
		km.Help,
		km.Quit,
		km.Return,
	}
}

var syncKeys = syncKeymap{
	Help: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("?"),
			k.WithHelp("?", "toggle help"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Auth)
			model.help.ShowAll = !model.help.ShowAll
			return model, nil
		},
	},
	Quit: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("q"),
			k.WithHelp("q", "quit"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			return m, tea.Quit
		},
	},
	Return: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("esc"),
			k.WithHelp("esc", "return"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			return m, func() tea.Msg { return state.HOME }
		},
	},
}
