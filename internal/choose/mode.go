// Package choose implements a generic options dialog. The caller supplies a title and a list
// of Option values. Each Option's Action is invoked with the shared Context when the user
// confirms the selection and returns the next Mode to transition to.
package choose

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"

	tea "charm.land/bubbletea/v2"
)

// Option is a single selectable item in the dialog. Label is shown in the list, Action is
// called when the user confirms this option.
type Option struct {
	Label  string
	Action func(ctx *core.Context) (core.Mode, tea.Cmd)
}

// Choose is a generic options dialog with cursor navigation.
type Choose struct {
	title   string
	options []Option
	cursor  int
	backTo  core.Mode
}

// New returns a Choose dialog with the given title and options. backTo is the mode to return
// to on cancel.
func New(title string, options []Option, backTo core.Mode) *Choose {
	return &Choose{title: title, options: options, backTo: backTo}
}

// SetCursor positions the highlight on the given row. Values outside the option range are
// silently ignored.
func (ch *Choose) SetCursor(i int) {
	if i >= 0 && i < len(ch.options) {
		ch.cursor = i
	}
}

// Init returns nil.
func (ch *Choose) Init(_ *core.Context) tea.Cmd { return nil }

// Update moves the cursor, confirms the selection, or cancels the dialog.
func (ch *Choose) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, chooseKeys.Up):
			if ch.cursor > 0 {
				ch.cursor--
			}
		case key.Matches(msg, chooseKeys.Down):
			if ch.cursor < len(ch.options)-1 {
				ch.cursor++
			}
		case key.Matches(msg, chooseKeys.Confirm):
			if ch.cursor >= 0 && ch.cursor < len(ch.options) {
				return ch.options[ch.cursor].Action(ctx)
			}
		case key.Matches(msg, chooseKeys.Cancel):
			if ch.backTo != nil {
				return ch.backTo, nil
			}
			return nil, nil
		}
	}
	return nil, nil
}

// Help returns the choose key bindings for the help footer.
func (ch *Choose) Help() help.KeyMap { return chooseKeys }

// NewPicker returns a Choose dialog pre configured from a Pickeable value's Options. When the
// user confirms a selection the value is written directly to ptr. after is called after the
// write, typically used for core.Rebuild.
func NewPicker[T config.Pickeable](title string, ptr *T, backTo core.Mode, after func(ctx *core.Context)) *Choose {
	options := (*ptr).Options()
	opts := make([]Option, len(options))
	for i, opt := range options {
		val := opt
		opts[i] = Option{
			Label: opt.PickLabel(),
			Action: func(ctx *core.Context) (core.Mode, tea.Cmd) {
				*ptr = val.(T)
				if after != nil {
					after(ctx)
				}
				return backTo, nil
			},
		}
	}
	ch := New(title, opts, backTo)
	current := config.Pickeable(*ptr)
	for i, opt := range options {
		if opt == current {
			ch.SetCursor(i)
			break
		}
	}
	return ch
}
