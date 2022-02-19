package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clstb/phi/go/pkg/ledger"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type unclassified struct {
	transactionIndex int
	postingIndex     int
	posting          ledger.Posting
}

type classifyModel struct {
	ledger         ledger.Ledger
	toClassify     []unclassified
	input          textinput.Model
	matches        []string
	matchSelection int
}

func newClassifyModel() classifyModel {
	input := textinput.New()
	input.Focus()
	return classifyModel{
		input: input,
	}
}

func (m classifyModel) Init() tea.Cmd { return nil }

func (m classifyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			s := msg.String()
			if s == "shift+tab" {
				m.matchSelection--
			} else {
				m.matchSelection++
			}

			if m.matchSelection > len(m.matches)-1 {
				m.matchSelection = 0
			} else if m.matchSelection < 0 {
				m.matchSelection = len(m.matches) - 1
			}
		case "enter":
			account := m.matches[m.matchSelection]

			var unclassified unclassified
			unclassified, m.toClassify = m.toClassify[0], m.toClassify[1:]
			unclassified.posting.Account = account

			transaction := m.ledger[unclassified.transactionIndex].(ledger.Transaction)
			transaction.Postings[unclassified.postingIndex] = unclassified.posting
			m.ledger[unclassified.transactionIndex] = transaction

			m.input.SetValue("")

			return m, func() tea.Msg {
				return m.ledger
			}
		default:
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

			return m, cmd
		}
	case ledger.Ledger:
		m.ledger = msg

		var toClassify []unclassified
		for i, v := range msg {
			transaction, ok := v.(ledger.Transaction)
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

	return m, nil
}

func (m classifyModel) View() string {
	if len(m.toClassify) == 0 {
		return "All classified"
	}

	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	renderOpens := func(opens []ledger.Open) string {
		sort.Slice(opens, func(i, j int) bool {
			return opens[i].Account < opens[j].Account
		})

		var view []string
		for _, open := range opens {
			if m.input.Value() == "" {
				view = append(view, open.Account)
				continue
			}

			s := blurredStyle.Render(open.Account)
			for i, match := range m.matches {
				if open.Account == match {
					s = open.Account
					if i == m.matchSelection {
						s = focusedStyle.Render(s)
					}
					break
				}
			}
			view = append(view, s)
		}
		return strings.Join(view, " \n")
	}

	var views []string
	for _, prefix := range []string{
		"Income",
		"Expenses",
		"Assets",
		"Equity",
		"Liabilities",
	} {
		views = append(views, renderOpens(m.ledger.Opens().Filter(func(open ledger.Open) bool {
			return strings.HasPrefix(open.Account, prefix)
		})))
	}

	s := lipgloss.JoinHorizontal(
		lipgloss.Top,
		views...,
	)

	transaction := m.ledger[m.toClassify[0].transactionIndex]

	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		transaction.String(),
		m.input.View(),
		s,
	)
}
