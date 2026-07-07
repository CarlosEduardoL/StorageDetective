// Package app is the top level Bubble Tea model. It owns the shared Context, holds the active
// Mode, and dispatches messages to whichever mode is current. The base view is drawn once per
// frame, the active mode's overlay composites on top.
//
// The modes are a full state machine. The current mode decides the next state by returning a
// Mode from Update, the app installs it and runs its Init. A sub mode returns to its caller
// by holding the caller's mode as a return target passed at construction.
//
// The app owns the long lived program level concerns, the modes own their real time
// behavior.
package app

import (
	"strings"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/explorer"
	"github.com/SolracHQ/stex/internal/styles"
	"github.com/SolracHQ/stex/internal/vfs"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// App is the top level Bubble Tea model. It owns the shared Context and the active Mode, and
// routes messages between them. Zero value is not valid, use New to construct one.
type App struct {
	ctx         *core.Context
	mode        core.Mode
	keys        Keys
	notifyQueue []core.AppNotify
}

// New constructs the top level Bubble Tea model with the given path, resolved config, and
// scanned root directory. It starts in the explorer mode and dispatches to whatever mode the
// user activates.
func New(path string, cfg config.Config, root *vfs.Dir) tea.Model {
	help := styles.HelpDefaults()

	table := table.New(
		table.WithFocused(true),
		table.WithStyles(styles.TableDefault()),
	)

	return &App{
		ctx: &core.Context{
			Path:    path,
			Config:  cfg,
			Width:   80,
			Height:  24,
			Root:    root,
			Current: root,
			Help:    help,
			Table:   table,
		},
		mode: &explorer.Explorer{},
		keys: DefaultKeys(),
	}
}

// Init returns the active mode's init command.
func (app *App) Init() tea.Cmd {
	if initCmd := app.mode.Init(app.ctx); initCmd != nil {
		return initCmd
	}
	return nil
}

// Update dispatches a message to the active mode. It also intercepts window resize (to keep
// the context in sync), the global quit key, notify messages, and the dismiss key. When
// a mode returns a new mode, the new mode's Init is called immediately and its command is
// appended.
func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		app.ctx.Width = sizeMsg.Width
		app.ctx.Height = sizeMsg.Height
		app.ctx.Help.SetWidth(sizeMsg.Width - 4)
		if app.ctx.Current != nil {
			core.Rebuild(app.ctx)
			core.UpdateInfo(app.ctx)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if key.Matches(msg, app.keys.Quit) {
			return app, tea.Quit
		}
		if len(app.notifyQueue) > 0 && key.Matches(msg, app.keys.DismissNotify) {
			return app, app.dismissNotify()
		}
		if key.Matches(msg, app.keys.HelpToggle) {
			app.ctx.Help.ShowAll = !app.ctx.Help.ShowAll
			return app, nil
		}
	case core.AppNotify:
		if !app.notifyAllowed(msg.Severity) {
			return app, nil
		}
		return app, app.enqueueNotify(msg)
	case core.AppNotifyDone:
		return app, app.dismissNotify()
	}

	next, cmd := app.mode.Update(app.ctx, msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if next != nil {
		app.mode = next
		if initCmd := next.Init(app.ctx); initCmd != nil {
			cmds = append(cmds, initCmd)
		}
	}

	return app, tea.Batch(cmds...)
}

// View returns the rendered frame. The base panels (title, table, info, footer) come from
// core.RenderBase, the active mode's overlay composites on top, and a notification appears at
// the top right.
func (app *App) View() tea.View {
	body := core.RenderBase(app.ctx, KeyMap{
		mode:    app.mode.Help(),
		globals: []key.Binding{app.keys.Quit, app.keys.HelpToggle},
	})
	if overlay := app.mode.Overlay(app.ctx); overlay != "" {
		body = overlayCenter(body, overlay)
	}
	if len(app.notifyQueue) > 0 {
		toast := RenderNotify(app.notifyQueue[0].Text, app.notifyQueue[0].Detail, app.notifyQueue[0].Severity)
		body = OverlayNotify(body, toast)
	}
	return core.WrapView(body)
}

// overlayCenter places foreground over background, centred inside the larger of the two. When
// foreground is empty the background is returned as is, when foreground fully covers background
// foreground is returned as is. ANSI escape sequences in the background are preserved by
// slicing with the ansi package instead of plain string ops.
func overlayCenter(background, foreground string) string {
	if foreground == "" || background == "" {
		return background
	}

	fgWidth, fgHeight := lipgloss.Size(foreground)
	bgWidth, bgHeight := lipgloss.Size(background)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return foreground
	}

	offsetX := (bgWidth - fgWidth) / 2
	offsetY := (bgHeight - fgHeight) / 2
	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}

	fgLines := strings.Split(foreground, "\n")
	bgLines := strings.Split(background, "\n")

	var buf strings.Builder
	for index, bgLine := range bgLines {
		if index > 0 {
			buf.WriteByte('\n')
		}
		if index < offsetY || index >= offsetY+fgHeight {
			buf.WriteString(bgLine)
			continue
		}

		pos := 0
		if offsetX > 0 {
			left := ansi.Truncate(bgLine, offsetX, "")
			pos = ansi.StringWidth(left)
			buf.WriteString(left)
			if pos < offsetX {
				buf.WriteString(strings.Repeat(" ", offsetX-pos))
				pos = offsetX
			}
		}

		fgLine := fgLines[index-offsetY]
		buf.WriteString(fgLine)
		pos += ansi.StringWidth(fgLine)

		right := ansi.TruncateLeft(bgLine, pos, "")
		bgW := ansi.StringWidth(bgLine)
		rightW := ansi.StringWidth(right)
		if rightW <= bgW-pos {
			buf.WriteString(strings.Repeat(" ", bgW-rightW-pos))
		}
		buf.WriteString(right)
	}

	return buf.String()
}
