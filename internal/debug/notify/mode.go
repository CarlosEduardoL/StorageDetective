// Package notify implements the notify form sub-mode for the debug menu.
package notify

import (
	"github.com/SolracHQ/stex/internal/choose"
	"github.com/SolracHQ/stex/internal/core"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Mode is the notify form sub-mode with three fields: title, detail, severity.
type Mode struct {
	returnTo core.Mode
	focus    int // 0=title, 1=detail, 2=severity
	title    textinput.Model
	detail   textinput.Model
	severity core.NotifySeverity
}

// New returns a notify form sub-mode that returns to returnTo on back.
func New(returnTo core.Mode) *Mode {
	title := textinput.New()
	title.Placeholder = "Short title"
	title.CharLimit = 100
	title.SetWidth(25)

	detail := textinput.New()
	detail.Placeholder = "Optional detail"
	detail.CharLimit = 200
	detail.SetWidth(25)

	return &Mode{
		returnTo: returnTo,
		focus:    0,
		title:    title,
		detail:   detail,
		severity: core.NotifyInfo,
	}
}

// Init focuses the title input so the user can start typing immediately.
func (m *Mode) Init(ctx *core.Context) tea.Cmd {
	return m.title.Focus()
}

// Update handles field navigation, Tab on severity, and firing the notification.
func (m *Mode) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Up):
			if m.focus > 0 {
				m.blur()
				m.focus--
				m.focusCurrent()
			}
			return nil, nil
		case key.Matches(keyMsg, keys.Down):
			if m.focus < 2 {
				m.blur()
				m.focus++
				m.focusCurrent()
			}
			return nil, nil
		case key.Matches(keyMsg, keys.Tab):
			if m.focus == 2 {
				options := []choose.Option{
					{Label: "Info", Action: func(ctx *core.Context) (core.Mode, tea.Cmd) {
						m.severity = core.NotifyInfo
						return m, nil
					}},
					{Label: "Warning", Action: func(ctx *core.Context) (core.Mode, tea.Cmd) {
						m.severity = core.NotifyWarn
						return m, nil
					}},
					{Label: "Error", Action: func(ctx *core.Context) (core.Mode, tea.Cmd) {
						m.severity = core.NotifyError
						return m, nil
					}},
				}
				ch := choose.New("Severity", options, m)
				ch.SetCursor(int(m.severity))
				return ch, nil
			}
			m.blur()
			m.focus++
			m.focusCurrent()
			return nil, nil
		case key.Matches(keyMsg, keys.Fire):
			cmd := core.NewNotifyCmd(m.title.Value(), m.detail.Value(), m.severity)
			return nil, cmd
		case key.Matches(keyMsg, keys.Back):
			return m.returnTo, nil
		}
	}

	if m.focus == 0 {
		var cmd tea.Cmd
		m.title, cmd = m.title.Update(msg)
		return nil, cmd
	}
	if m.focus == 1 {
		var cmd tea.Cmd
		m.detail, cmd = m.detail.Update(msg)
		return nil, cmd
	}
	return nil, nil
}

func (m *Mode) Help() help.KeyMap { return keys }
func (m *Mode) Name() string      { return "notify" }

func (m *Mode) blur() {
	switch m.focus {
	case 0:
		m.title.Blur()
	case 1:
		m.detail.Blur()
	}
}

func (m *Mode) focusCurrent() {
	switch m.focus {
	case 0:
		m.title.Focus()
	case 1:
		m.detail.Focus()
	}
}
