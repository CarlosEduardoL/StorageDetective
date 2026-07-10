package app

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/styles"
	"github.com/charmbracelet/x/ansi"
)

const (
	leftGlyphChar  = ""
	rightGlyphChar = ""
)

// bgColor returns the background color constant for a powerbar segment.
func bgColor(idx int) string {
	if idx == 0 {
		return styles.AccentColor
	}
	if idx%2 == 1 {
		return styles.DimColor
	}
	return styles.DimAltColor
}

// renderLabel styles a segment label with the background color at index.
func renderLabel(label string, idx int) string {
	switch {
	case idx == 0:
		return styles.PowerbarMode.Render(label)
	case idx%2 == 1:
		return styles.PowerbarSegment.Render(label)
	default:
		return styles.PowerbarSegmentAlt.Render(label)
	}
}

// pad wraps a string in a space on each side.
func pad(s string) string { return " " + s + " " }

// powerbar renders the application status line.
func (app *App) powerbar() string {
	width := app.screen.Width
	showGlyphs := app.ctx.Config.ShowPowerGlyphs
	icons := app.ctx.Config.ShowIcons
	current := app.ctx.Current
	cfg := app.ctx.Config

	var leftGlyph, rightGlyph string
	if showGlyphs {
		leftGlyph = leftGlyphChar
		rightGlyph = rightGlyphChar
	}

	var left []string
	left = append(left, pad(app.mode.Name()))
	left = append(left, pad(config.GroupingString(cfg.Grouping)))
	if cfg.Filter != nil {
		pattern := cfg.Filter.String()
		if icons {
			left = append(left, " 🔍 /"+pattern+"/ ")
		} else {
			left = append(left, " /"+pattern+"/ ")
		}
	}
	left = append(left, pad(hiddenIndicator(icons, cfg.ShowHidden)))

	right := []string{
		pad(current.Size().String()),
		pad(fileCountStr(icons, current.Count().Files)),
		pad(dirCountStr(icons, current.Count().Dirs)),
	}

	var builder strings.Builder
	for i, l := range left {
		if i > 0 {
			builder.WriteString(styles.SepString(leftGlyph, bgColor(i-1), bgColor(i)))
		}
		builder.WriteString(renderLabel(l, i))
	}
	if showGlyphs && len(left) > 0 {
		builder.WriteString(styles.SepString(leftGlyph, bgColor(len(left)-1), ""))
	}
	leftRendered := builder.String()
	leftWidth := ansi.StringWidth(leftRendered)

	builder.Reset()
	fixedCount := len(right)
	if showGlyphs && fixedCount > 0 {
		builder.WriteString(styles.SepString(rightGlyph, bgColor(fixedCount-1), ""))
	}
	for i, l := range right {
		if i > 0 {
			builder.WriteString(styles.SepString(rightGlyph, bgColor(fixedCount-1-i), bgColor(fixedCount-i)))
		}
		builder.WriteString(renderLabel(l, fixedCount-1-i))
	}
	rightFixedLabel := builder.String()
	rightFixedWidth := ansi.StringWidth(rightFixedLabel)

	available := width - leftWidth - rightFixedWidth - 3
	if available > 0 {
		right = append(right, pad(shortenTitlePath(current.FullPath(), available)))
	}

	builder.Reset()
	rightCount := len(right)
	if showGlyphs && rightCount > 0 {
		builder.WriteString(styles.SepString(rightGlyph, bgColor(rightCount-1), ""))
	}
	for i, l := range right {
		if i > 0 {
			builder.WriteString(styles.SepString(rightGlyph, bgColor(rightCount-1-i), bgColor(rightCount-i)))
		}
		builder.WriteString(renderLabel(l, rightCount-1-i))
	}
	rightLabel := builder.String()
	rightWidth := ansi.StringWidth(rightLabel)

	padding := max(width-leftWidth-rightWidth, 0)
	return leftRendered + strings.Repeat(" ", padding) + rightLabel
}

// fileCountStr returns a formatted file count, optionally with an icon.
func fileCountStr(icons bool, files int) string {
	if icons {
		return fmt.Sprintf("%d 📄", files)
	}
	return fmt.Sprintf("%d f", files)
}

// dirCountStr returns a formatted directory count, optionally with an icon.
func dirCountStr(icons bool, dirs int) string {
	if icons {
		return fmt.Sprintf("%d 📁", dirs)
	}
	return fmt.Sprintf("%d d", dirs)
}

// hiddenIndicator returns a short label showing whether hidden files are visible.
func hiddenIndicator(icons, shown bool) string {
	if icons {
		if shown {
			return "👁"
		}
		return "🙈"
	}
	if shown {
		return "sh"
	}
	return "h"
}

// shortenTitlePath truncates a path with a /.../ middle ellipsis when it exceeds maxLen.
func shortenTitlePath(path string, maxLen int) string {
	if maxLen < 1 || path == "" {
		return ""
	}
	if len(path) <= maxLen {
		return path
	}
	if maxLen < 5 {
		return path[:maxLen]
	}
	keep := maxLen - 5
	return "/.../" + path[len(path)-keep:]
}
