// Package core is the app's base architecture. The base view is drawn once per frame, modes
// contribute a keymap and can return an overlay that composites on top. The model is the nvim
// style, a stable base, modes that change keys, overlays that add context, modes that compose.
//
// The app is the long lived piece, it owns the Context and any program lifetime orchestration.
// Modes are transient states for real time behavior. A sync that needs to lock input and mouse
// is the exception, it lives in a mode so the lock is enforced by the mode being active.
package core

import (
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/layout"
	"github.com/SolracHQ/stex/internal/vfs"

	"charm.land/bubbles/v2/table"
)

// Context is the mutable state bag passed to every mode. Modes read and write the fields
// directly. The scanned tree (Root, Current) lives here so a mode transition does not need to
// rebuild it.
type Context struct {
	screen layout.Rect // set by SetScreen, read via Screen()
	Path   string

	Root, Current *vfs.Dir
	Config        config.Config

	// Table widget, sized by SetScreen. Modes may scroll or resize it.
	Table table.Model
	Info  InfoState // cached right pane content
	Items []vfs.FileSystemItem
}

// InfoState holds the cached right pane stat result and current layout bounds.
type InfoState struct {
	Path   string // cursor path, for stat cache
	file   *FileInfo
	Bounds layout.Rect // info panel bounds, set by SetScreen
}

// Screen returns the content area rectangle.
func (ctx *Context) Screen() layout.Rect { return ctx.screen }

// SetScreen updates the content area, table dimensions, and info pane layout.
func (ctx *Context) SetScreen(s layout.Rect) {
	ctx.screen = s
	if s.Width >= SplitViewThreshold {
		ctx.Info.Bounds = s.ShrinkSides(s.Width/2, 0, 0, 0)
	} else {
		ctx.Info.Bounds = layout.Rect{}
	}
	if s.Width > 0 && s.Height > 0 {
		ctx.Table.SetWidth(s.Width)
		ctx.Table.SetHeight(s.Height)
	}
}
