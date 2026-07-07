package notify

import "charm.land/bubbles/v2/key"

type formKeys struct {
	Up, Down, Tab, Fire, Back key.Binding
}

var keys = formKeys{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "previous field"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next field"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle severity"),
	),
	Fire: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "fire notify"),
	),
	Back: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "back to menu"),
	),
}

func (k formKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Tab, k.Fire, k.Back}
}

func (k formKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Tab, k.Fire, k.Back}}
}
