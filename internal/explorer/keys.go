package explorer

import "charm.land/bubbles/v2/key"

type keys struct {
	Up, Down, Enter, Back key.Binding
	Search, ClearFilter   key.Binding
	Command, Settings     key.Binding
}

// explorerKeys is the singleton key map used by the explorer. The list is split into four rows
// for the expanded help view.
var explorerKeys = keys{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j/s", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "l", "right", "d"),
		key.WithHelp("→/l/d", "move in"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace", "left", "h", "a"),
		key.WithHelp("←/h/a", "move out"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter by name"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear filter"),
	),
	Command: key.NewBinding(
		key.WithKeys(":"),
		key.WithHelp(":", "command"),
	),
	Settings: key.NewBinding(
		key.WithKeys(","),
		key.WithHelp(",", "settings"),
	),
}

// ShortHelp returns the bindings shown in the collapsed footer.
func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back}
}

// FullHelp returns the bindings grouped into rows for the expanded help view.
func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Back},
		{k.Search, k.ClearFilter, k.Command, k.Settings},
	}
}
