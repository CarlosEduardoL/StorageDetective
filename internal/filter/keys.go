package filter

import "charm.land/bubbles/v2/key"

// keys holds the filter mode's key bindings.
type keys struct {
	Confirm    key.Binding
	Cancel     key.Binding
	ToggleLive key.Binding
}

func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Cancel, k.ToggleLive}
}

func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Confirm, k.Cancel, k.ToggleLive}}
}

// filterKeys is the singleton key map used in filter mode.
var filterKeys = keys{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "keep filter"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	ToggleLive: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "toggle live filter"),
	),
}
