package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	k "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/internal/phi/commands"
	"github.com/clstb/phi/go/internal/phi/key"
	"github.com/clstb/phi/go/internal/phi/state"
	"github.com/clstb/phi/go/pkg/client"
)

type Home struct {
	client *client.Client
	help   help.Model
}

func NewHome(client *client.Client) Home {
	return Home{
		client: client,
		help:   help.New(),
	}
}

func (m Home) Init() tea.Cmd {
	return nil
}

func (m Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		for _, key := range homeKeys.slice() {
			if key.Match(msg.String()) {
				newModel, cmd := key.Action(m)
				m, cmds = newModel.(Home), append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Home) View() string {
	return fmt.Sprintf(
		"Welcome to phi!\n\n%s",
		m.help.View(homeKeys),
	)
}

type homeKeymap struct {
	Help     key.Key
	Quit     key.Key
	Sync     key.Key
	Classify key.Key
	Link     key.Key
}

func (km homeKeymap) ShortHelp() []k.Binding {
	return []k.Binding{km.Help.Binding, km.Quit.Binding}
}

func (km homeKeymap) FullHelp() [][]k.Binding {
	return [][]k.Binding{
		{km.Sync.Binding, km.Classify.Binding, km.Link.Binding},
		{km.Help.Binding, km.Quit.Binding},
	}
}

func (km homeKeymap) slice() []key.Key {
	return []key.Key{
		km.Help,
		km.Quit,
		km.Sync,
		km.Classify,
		km.Link,
	}
}

var homeKeys = homeKeymap{
	Help: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("?"),
			k.WithHelp("?", "toggle help"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Home)
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
	Sync: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("s"),
			k.WithHelp("s", "synchronize ledger"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			return m, func() tea.Msg { return state.SYNC }
		},
	},
	Classify: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("c"),
			k.WithHelp("c", "classify transactions"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			return m, func() tea.Msg { return state.CLASSIFY }
		},
	},
	Link: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("l"),
			k.WithHelp("l", "link bank account"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Home)
			return m, commands.OpenLink(model.client)
		},
	},
}
