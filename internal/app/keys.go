package app

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

// Keys holds the key bindings that the app intercepts before forwarding to the active mode.
type Keys struct {
	Quit          key.Binding
	DismissNotify key.Binding
	HelpToggle    key.Binding
}

// DefaultKeys returns the standard app level key bindings.
func DefaultKeys() Keys {
	return Keys{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+d"),
			key.WithHelp("ctrl+c", "quit"),
		),
		DismissNotify: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "dismiss"),
		),
		HelpToggle: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "help"),
		),
	}
}

// KeyMap wraps a mode's help.KeyMap with the app-level global bindings so they appear in the
// help footer alongside the mode's own keys.
type KeyMap struct {
	mode    help.KeyMap
	globals []key.Binding
}

func (km KeyMap) ShortHelp() []key.Binding {
	return append(km.mode.ShortHelp(), km.globals...)
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return append(km.mode.FullHelp(), km.globals)
}
