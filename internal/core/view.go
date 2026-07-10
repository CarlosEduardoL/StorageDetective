package core

import (
	"strings"

	"github.com/SolracHQ/stex/internal/layout"
	"github.com/SolracHQ/stex/internal/styles"

	"charm.land/lipgloss/v2"
)

// SplitViewThreshold is the minimum inner width at which the layout switches from a single
// column to the table plus info side by side.
const SplitViewThreshold = 80

// RenderBase renders the table and info panels into the full usable area.
func RenderBase(ctx *Context) string {
	if ctx.Screen().Width == 0 || ctx.Screen().Height == 0 || ctx.Current == nil {
		return ""
	}

	if ctx.Screen().Width >= SplitViewThreshold {
		return splitView(ctx, ctx.Screen())
	}
	return narrowView(ctx, ctx.Screen())
}

// splitView renders the table and info panels side by side.
func splitView(ctx *Context, inner layout.Rect) string {
	infoRect := ctx.Info.Bounds
	leftWidth := inner.Width - infoRect.Width - 1

	left := ctx.Table.View()
	right := InfoContent(ctx)
	if right == "" {
		return left
	}

	lines := max(strings.Count(left, "\n"), strings.Count(right, "\n")) + 1
	sep := styles.Dim.Render(strings.Repeat("│\n", lines-1) + "│")
	leftBlock := lipgloss.NewStyle().Width(leftWidth).Render(left)
	rightBlock := lipgloss.NewStyle().Width(infoRect.Width).Render(right)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftBlock, sep, rightBlock)
}

// narrowView renders the single column layout used on terminals narrower than the split
// threshold. The info pane is hidden.
func narrowView(ctx *Context, _ layout.Rect) string {
	return ctx.Table.View()
}
