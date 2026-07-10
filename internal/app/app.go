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
	"github.com/SolracHQ/stex/internal/layout"
	"github.com/SolracHQ/stex/internal/styles"
	"github.com/SolracHQ/stex/internal/vfs"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// App is the top level Bubble Tea model. Routes messages between modes. Zero value is not
// valid, use New.
type App struct {
	ctx    *core.Context
	mode   core.Mode
	help   help.Model // help widget, rendered by View at the top
	screen layout.Rect  // full terminal rect, set on resize
	keys   Keys

	notifyQueue []core.AppNotify
}

// New constructs the top level Bubble Tea model with the given path, resolved config, and
// scanned root directory. It starts in the explorer mode and dispatches to whatever mode the
// user activates.
func New(path string, cfg config.Config, root *vfs.Dir) tea.Model {
	screen := layout.New(0, 0, 80, 24)

	tbl := table.New(
		table.WithFocused(true),
		table.WithStyles(styles.TableDefault()),
	)

	ctx := &core.Context{
		Path:    path,
		Config:  cfg,
		Root:    root,
		Current: root,
		Table:   tbl,
	}
	ctx.SetScreen(screen.Shrink(1))
	return &App{
		ctx:    ctx,
		mode:   &explorer.Explorer{},
		help:   styles.HelpDefaults(),
		screen: screen,
		keys:   DefaultKeys(),
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
		app.screen = layout.New(0, 0, sizeMsg.Width, sizeMsg.Height)
		app.help.SetWidth(sizeMsg.Width - 4)
		app.computeScreen()
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
			app.help.ShowAll = !app.help.ShowAll
			app.computeScreen()
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
		app.computeScreen()
		if initCmd := next.Init(app.ctx); initCmd != nil {
			cmds = append(cmds, initCmd)
		}
	}

	return app, tea.Batch(cmds...)
}

func (app *App) km() KeyMap {
	return KeyMap{
		mode:    app.mode.Help(),
		globals: []key.Binding{app.keys.Quit, app.keys.HelpToggle},
	}
}

// computeScreen updates ctx.Screen based on the current help size.
func (app *App) computeScreen() {
	inner := app.screen.Shrink(1)

	renderedHelp := app.help.View(app.km())
	helpHeight := 0
	if renderedHelp != "" {
		centered := styles.CenterText(renderedHelp, inner.Width)
		helpHeight = strings.Count(centered, "\n") + 1
	}

	_, rest := inner.Top(helpHeight)
	content, _ := rest.Bottom(1)
	app.ctx.SetScreen(content)
}

// View renders the application layout.
func (app *App) View() tea.View {
	inner := app.screen.Shrink(1)

	renderedHelp := app.help.View(app.km())
	var centeredHelp string
	helpHeight := 0
	if renderedHelp != "" {
		centeredHelp = styles.CenterText(renderedHelp, inner.Width)
		helpHeight = strings.Count(centeredHelp, "\n") + 1
	}

	app.ctx.Table.SetHeight(app.ctx.Screen().Height)
	body := core.RenderBase(app.ctx)
	bordered := styles.BorderNorm.Render(body)

	powerbar := app.powerbar()

	var frame string
	if helpHeight > 0 {
		frame = centeredHelp + "\n"
	}
	frame += bordered + "\n" + powerbar

	if overlay := app.mode.Overlay(app.ctx); overlay != "" {
		frame = overlayCenter(frame, overlay)
	}
	if len(app.notifyQueue) > 0 {
		toast := renderNotify(app.notifyQueue[0].Text, app.notifyQueue[0].Detail, app.notifyQueue[0].Severity)
		frame = overlayNotify(frame, toast)
	}

	view := tea.NewView(frame)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

// overlayCenter composites foreground over background, centred.
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
