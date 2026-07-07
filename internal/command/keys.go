package command

import "charm.land/bubbles/v2/key"

type keys struct {
	Confirm          key.Binding
	Cancel           key.Binding
	AcceptSuggestion key.Binding
	NextSuggestion   key.Binding
	PrevSuggestion   key.Binding
}

func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.AcceptSuggestion, k.NextSuggestion, k.PrevSuggestion, k.Cancel}
}

func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Confirm, k.AcceptSuggestion, k.NextSuggestion, k.PrevSuggestion, k.Cancel}}
}

var commandKeys = keys{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "run"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	AcceptSuggestion: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "accept completion"),
	),
	NextSuggestion: key.NewBinding(
		key.WithKeys("down", "shift+tab"),
		key.WithHelp("↓/tab", "next suggestion"),
	),
	PrevSuggestion: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "prev suggestion"),
	),
}
