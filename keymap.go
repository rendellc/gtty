package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	CycleView key.Binding
	Clear key.Binding
	Quit      key.Binding
}


func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down}, 
		{k.Clear},
		{k.CycleView, k.Quit}, 
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	CycleView: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "cycle views"),
	),
	Clear: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("delete", "clear terminal"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}
