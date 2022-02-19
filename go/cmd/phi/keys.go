package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Help       key.Binding
	Quit       key.Binding
	Sync       key.Binding
	Classify   key.Binding
	AddAccount key.Binding
	Link       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Sync, k.Classify, k.AddAccount, k.Link},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Sync: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "synchronize transactions and accounts"),
	),
	Classify: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "classify unassigned transactions"),
	),
	AddAccount: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "add account"),
	),
	Link: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "link bank account"),
	),
}
