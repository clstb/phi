package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/internal/phi/beancount"
	"github.com/clstb/phi/go/internal/phi/commands"
	"github.com/clstb/phi/go/internal/phi/models"
	"github.com/clstb/phi/go/internal/phi/state"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/urfave/cli/v2"
)

type model struct {
	ctx       *cli.Context
	state     state.State
	subModels map[state.State]tea.Model
	client    *client.Client
}

func newModel(ctx *cli.Context) model {
	client := client.NewClient(ctx.String("api-url"))

	return model{
		ctx:   ctx,
		state: state.AUTH,
		subModels: map[state.State]tea.Model{
			state.HOME:     models.NewHome(client),
			state.AUTH:     models.NewAuth(client),
			state.CLASSIFY: models.NewClassify(),
			state.SYNC:     models.NewSync(client),
		},
		client: client,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		commands.LoadLedger(m.ctx.Path("ledger")),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case state.State:
		m.state = msg
		cmds = append(cmds, m.subModels[msg].Init())
	case tea.KeyMsg:
		var cmd tea.Cmd
		m.subModels[m.state], cmd = m.subModels[m.state].Update(msg)
		cmds = append(cmds, cmd)
	case error:
		panic(msg)
	default:
		switch msg := msg.(type) {
		case client.Session:
			m.client.SetBearerToken(msg.Token)
		case beancount.Ledger:
			cmds = append(cmds, commands.SaveLedger(m.ctx.Path("ledger"), msg))
		}
		for state, subModel := range m.subModels {
			var cmd tea.Cmd
			m.subModels[state], cmd = subModel.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.subModels[m.state].View()
}
