package core

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

// Tick returns a tea.Cmd that sends the given message after the specified duration. A wrapper
// around tea.Tick that avoids the inline closure boilerplate.
func Tick(d time.Duration, msg tea.Msg) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return msg
	})
}
