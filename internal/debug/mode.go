// Package debug implements a diagnostic mode with a cursor-driven action menu. Actions
// transition to sub-modes under the same internal/debug/ tree.
package debug

import (
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/debug/notify"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Mode is the debug menu mode. Zero value is valid.
type Mode struct {
	returnTo core.Mode
	cursor   int
}

// New returns a debug menu mode that returns to returnTo on back.
func New(returnTo core.Mode) *Mode {
	return &Mode{returnTo: returnTo}
}

func (Mode) Init(*core.Context) tea.Cmd { return nil }

// Update navigates the menu, selects actions, or returns to the caller.
func (m *Mode) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(keyMsg, keys.Down):
			if m.cursor < len(menuActions)-1 {
				m.cursor++
			}
		case key.Matches(keyMsg, keys.Select):
			switch m.cursor {
			case 0:
				return notify.New(m), nil
			}
		case key.Matches(keyMsg, keys.Back):
			return m.returnTo, nil
		}
	}
	return nil, nil
}

func (m *Mode) Help() help.KeyMap { return keys }
func (m *Mode) Name() string      { return "debug" }
