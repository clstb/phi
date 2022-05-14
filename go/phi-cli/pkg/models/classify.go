package models

import (
	"fmt"
	beancount2 "github.com/clstb/phi/go/phi-cli/pkg/beancount"
	"github.com/clstb/phi/go/phi-cli/pkg/key"
	"github.com/clstb/phi/go/phi-cli/pkg/state"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	k "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type unclassified struct {
	transactionIndex int
	postingIndex     int
	posting          beancount2.Posting
}

type Classify struct {
	ledger         beancount2.Ledger
	toClassify     []unclassified
	matches        []string
	matchSelection int
	input          textinput.Model
	help           help.Model
}

func NewClassify() Classify {
	input := textinput.New()
	input.Focus()
	return Classify{
		input: input,
		help:  help.New(),
	}
}

func (m Classify) Init() tea.Cmd { return nil }

func (m Classify) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		for _, key := range classifyKeys.slice() {
			if key.Match(msg.String()) {
				newModel, cmd := key.Action(m)
				m, cmds = newModel.(Classify), append(cmds, cmd)
			}
		}
		switch msg.String() {
		case "?", "tab", "shift+tab":
		default:
			cmds = append(cmds, m.updateInputs(msg))
		}
	case beancount2.Ledger:
		m.ledger = msg

		var toClassify []unclassified
		for i, v := range msg {
			transaction, ok := v.(beancount2.Transaction)
			if !ok {
				continue
			}
			for j, posting := range transaction.Postings {
				if strings.HasSuffix(posting.Account, "Unassigned") {
					toClassify = append(toClassify, unclassified{
						transactionIndex: i,
						postingIndex:     j,
						posting:          posting,
					})
				}
			}
		}
		m.toClassify = toClassify
	}

	return m, tea.Batch(cmds...)
}

func (m Classify) View() string {
	if len(m.toClassify) == 0 {
		return "All classified"
	}

	transaction := m.ledger[m.toClassify[0].transactionIndex]

	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	var groupedAccounts []string
	opens := m.ledger.Opens()
	sort.Slice(opens, func(i, j int) bool { return opens[i].Account < opens[j].Account })
	for _, prefix := range []string{
		"Income",
		"Expenses",
		"Assets",
		"Equity",
		"Liabilities",
	} {
		var accounts []string
		for _, account := range opens.Filter(func(open beancount2.Open) bool {
			return strings.HasPrefix(open.Account, prefix)
		}) {
			if m.input.Value() == "" {
				accounts = append(accounts, account.Account)
				continue
			}

			s := blurredStyle.Render(account.Account)
			for i, match := range m.matches {
				if account.Account == match {
					s = account.Account
					if i == m.matchSelection {
						s = focusedStyle.Render(s)
					}
					break
				}
			}
			accounts = append(accounts, s)
		}
		groupedAccounts = append(groupedAccounts, strings.Join(accounts, " \n"))
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s\n\n%s",
		transaction.String(),
		m.input.View(),
		lipgloss.JoinHorizontal(lipgloss.Top, groupedAccounts...),
		m.help.View(classifyKeys),
	)
}

func (m *Classify) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	var matches []string
	for _, open := range m.ledger.Opens() {
		if fuzzy.MatchFold(m.input.Value(), open.Account) {
			matches = append(matches, open.Account)
		}
	}
	m.matches = matches
	m.matchSelection = 0

	return cmd
}

type classifyKeymap struct {
	Help           key.Key
	Return         key.Key
	SelectNext     key.Key
	SelectPrevious key.Key
	Submit         key.Key
}

func (km classifyKeymap) ShortHelp() []k.Binding {
	return []k.Binding{km.Help.Binding, km.Return.Binding}
}

func (km classifyKeymap) FullHelp() [][]k.Binding {
	return [][]k.Binding{
		{km.SelectNext.Binding, km.SelectPrevious.Binding, km.Submit.Binding},
		{km.Return.Binding, km.Help.Binding},
	}
}

func (km classifyKeymap) slice() []key.Key {
	return []key.Key{
		km.Help,
		km.Return,
		km.SelectNext,
		km.SelectPrevious,
		km.Submit,
	}
}

var classifyKeys = classifyKeymap{
	Help: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("?"),
			k.WithHelp("?", "toggle help"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Classify)
			model.help.ShowAll = !model.help.ShowAll
			return model, nil
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
	SelectNext: key.Key{
		Binding: k.NewBinding(
			k.WithKeys("tab"),
			k.WithHelp("tab", "select next element"),
		),
		Action: func(m tea.Model) (tea.Model, tea.Cmd) {
			model := m.(Classify)
			model.matchSelection++
			if model.matchSelection > len(model.matches)-1 {
				model.matchSelection = 0
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
			model := m.(Classify)
			model.matchSelection--
			if model.matchSelection < 0 {
				model.matchSelection = len(model.matches) - 1
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
			model := m.(Classify)

			if len(model.matches) == 0 {
				model.ledger = append(model.ledger, beancount2.Open{
					Date:    "1970-01-01",
					Account: model.input.Value(),
				})
			} else {
				account := model.matches[model.matchSelection]

				var unclassified unclassified
				unclassified, model.toClassify = model.toClassify[0], model.toClassify[1:]
				unclassified.posting.Account = account

				transaction := model.ledger[unclassified.transactionIndex].(beancount2.Transaction)
				transaction.Postings[unclassified.postingIndex] = unclassified.posting
				model.ledger[unclassified.transactionIndex] = transaction
			}

			model.input.Reset()
			return model, func() tea.Msg { return model.ledger }
		},
	},
}
