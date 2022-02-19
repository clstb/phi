package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/pkg/tink"
)

type banksModel struct {
	tinkClient *tink.Client
	consents   []tink.ProviderConsent
}

func newBanksModel() banksModel {
	return banksModel{}
}

func (m banksModel) Init() tea.Cmd {
	return func() tea.Msg {
		consents, err := m.tinkClient.ProviderConsents()
		if err != nil {
			return err
		}
		return consents
	}
}

func (m banksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *tink.Client:
		m.tinkClient = msg
	case []tink.ProviderConsent:
		m.consents = msg
	}

	return m, nil
}

func (m banksModel) View() string {
	var b strings.Builder
	for _, consent := range m.consents {
		b.WriteString(fmt.Sprintf("%s %s\n", consent.ProviderName, consent.Status))
	}
	return b.String()
}
