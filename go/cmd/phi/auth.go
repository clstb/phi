package main

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clstb/phi/go/pkg/auth/pb"
	authpb "github.com/clstb/phi/go/pkg/auth/pb"
	"github.com/clstb/phi/go/pkg/config"
)

type authModel struct {
	ctx        context.Context
	authClient authpb.AuthClient

	focusIndex int
	inputs     []textinput.Model
}

func newAuthModel(ctx context.Context) authModel {
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

	return authModel{
		ctx:    ctx,
		inputs: inputs,
	}
}

func (m authModel) Init() tea.Cmd { return nil }

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case authpb.AuthClient:
		m.authClient = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, login(
					m.ctx,
					m.authClient,
					&pb.User{
						Name:     m.inputs[0].Value(),
						Password: m.inputs[1].Value(),
					},
				)
			}
			if s == "enter" && m.focusIndex == len(m.inputs)+1 {
				return m, register(
					m.ctx,
					m.authClient,
					&pb.User{
						Name:     m.inputs[0].Value(),
						Password: m.inputs[1].Value(),
					},
				)
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)+1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) + 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	return m, m.updateInputs(msg)
}

func (m *authModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m authModel) View() string {
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
		b.WriteString(focusedStyle.Render("\n\tLogin"))
		b.WriteString(blurredStyle.Render("\tRegister"))
	} else if m.focusIndex == len(m.inputs)+1 {
		b.WriteString(blurredStyle.Render("\n\tLogin"))
		b.WriteString(focusedStyle.Render("\tRegister"))
	} else {
		b.WriteString(blurredStyle.Render("\n\tLogin\tRegister"))
	}

	return b.String()
}

func login(
	ctx context.Context,
	authClient authpb.AuthClient,
	user *pb.User,
) tea.Cmd {
	return func() tea.Msg {
		token, err := authClient.Login(ctx, user)
		if err != nil {
			return err
		}
		return config.PhiToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			ExpiresAt:    token.ExpiresAt,
			Scope:        token.Scope,
		}
	}
}

func register(
	ctx context.Context,
	authClient authpb.AuthClient,
	user *pb.User,
) tea.Cmd {
	return func() tea.Msg {
		token, err := authClient.Register(ctx, user)
		if err != nil {
			return err
		}
		return config.PhiToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			ExpiresAt:    token.ExpiresAt,
			Scope:        token.Scope,
		}
	}
}
