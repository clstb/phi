package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	k "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clstb/phi/go/internal/phi/key"
	"github.com/clstb/phi/go/internal/phi/state"
	"github.com/clstb/phi/go/pkg/client"
)

type Auth struct {
	focusIndex int
	inputs     []textinput.Model
	help       help.Model
	client     *client.Client
}

func NewAuth(client *client.Client) Auth {
	inputs := make([]textinput.Model, 2)

	for i := range inputs {
		t := textinput.New()

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		inputs[i] = t
	}

	return Auth{
		inputs: inputs,
		help:   help.NewModel(),
		client: client,
	}
}

func (m Auth) Init() tea.Cmd { return nil }

func (m Auth) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		for _, key := range authKeys.slice() {
			if key.Match(msg.String()) {
				newModel, cmd := key.Action(m)
				m, cmds = newModel.(Auth), append(cmds, cmd)
			}
		}
		if msg.String() != "?" {
			cmds = append(cmds, m.updateInputs(msg))
		}
	case client.Session:
		cmds = append(cmds, func() tea.Msg { return state.HOME })
	}

	return m, tea.Batch(cmds...)
}

func (m Auth) View() string {
	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	if m.focusIndex == len(m.inputs) {
		b.WriteString(focusedStyle.Render("\n\n\tLogin"))
		b.WriteString(blurredStyle.Render("\tRegister"))
	} else if m.focusIndex == len(m.inputs)+1 {
		b.WriteString(blurredStyle.Render("\n\n\tLogin"))
		b.WriteString(focusedStyle.Render("\tRegister"))
	} else {
		b.WriteString(blurredStyle.Render("\n\n\tLogin\tRegister"))
	}

	return fmt.Sprintf("%s\n\n%s", b.String(), m.help.View(authKeys))
}

func (m *Auth) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *Auth) updateFocus() {
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focusIndex {
			m.inputs[i].Focus()
			continue
		}
		m.inputs[i].Blur()
	}
}

type authKeymap struct {
	Help           key.Key
	Quit           key.Key
	SelectNext     key.Key
	SelectPrevious key.Key
	Submit         key.Key
}

func (km authKeymap) ShortHelp() []k.Binding {
	return []k.Binding{km.Help.Binding, km.Quit.Binding}
}

func (km authKeymap) FullHelp() [][]k.Binding {
	return [][]k.Binding{
		{km.SelectNext.Binding, km.SelectPrevious.Binding, km.Submit.Binding},
		{km.Quit.Binding, km.Help.Binding},
	}
}

func (km authKeymap) slice() []key.Key {
	return []key.Key{
		km.Help,
		km.Quit,
		km.SelectNext,
		km.SelectPrevious,
		km.Submit,
	}
}

var authKeys = authKeymap{
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
			k.WithKeys("esc"),
			k.WithHelp("esc", "quit"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			return m, tea.Quit
		},
	},
	SelectNext: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("tab"),
			k.WithHelp("tab", "select next element"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Auth)
			defer model.updateFocus()
			model.focusIndex++
			if model.focusIndex > len(model.inputs)+1 {
				model.focusIndex = 0
			}
			return model, nil
		},
	},
	SelectPrevious: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("shift+tab"),
			k.WithHelp("shift+tab", "select previous element"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Auth)
			defer model.updateFocus()
			model.focusIndex--
			if model.focusIndex < 0 {
				model.focusIndex = len(model.inputs) + 1
			}
			return model, nil
		},
	},
	Submit: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("enter"),
			k.WithHelp("enter", "submit"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Auth)
			switch {
			case model.focusIndex == len(model.inputs):
				return model, func() tea.Msg {
					session, err := model.client.Login(
						model.inputs[0].Value(),
						model.inputs[1].Value(),
					)
					if err != nil {
						return err
					}
					return session
				}
			case model.focusIndex == len(model.inputs)+1:
				return model, func() tea.Msg {
					session, err := model.client.Register(
						model.inputs[0].Value(),
						model.inputs[1].Value(),
					)
					if err != nil {
						return err
					}
					return session
				}
			default:
				defer model.updateFocus()
				model.focusIndex++
				if model.focusIndex > len(model.inputs)+1 {
					model.focusIndex = 0
				}
				return model, nil
			}
		},
	},
}
