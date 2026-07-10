package core

import (
	"fmt"
	"image/color"
	"math"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/vfs"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

// Column widths used when building the table columns.
const (
	sizePctWidth = 9
	sizeWidth    = 12
)

// buildColumns sets the table column headers and widths based on the sort direction.
func buildColumns(ctx *Context) {
	nameWidth := max(ctx.Screen().Width-sizePctWidth-sizeWidth, 1)

	var sizeLabel, nameLabel string
	switch ctx.Config.SortBy {
	case config.SortBySize:
		if ctx.Config.SortOrder == config.Descending {
			sizeLabel = "Size↓"
		} else {
			sizeLabel = "Size↑"
		}
		nameLabel = "Name"
	case config.SortByName:
		sizeLabel = "Size"
		if ctx.Config.SortOrder == config.Descending {
			nameLabel = "Name↓"
		} else {
			nameLabel = "Name↑"
		}
	}

	ctx.Table.SetColumns([]table.Column{
		{Title: " Size%", Width: sizePctWidth},
		{Title: sizeLabel, Width: sizeWidth},
		{Title: nameLabel, Width: nameWidth},
	})
}

// buildRows transforms items into bubbletea table rows and sets them on the table.
func buildRows(ctx *Context, items []vfs.FileSystemItem) {
	rows := make([]table.Row, 0, len(items))
	for _, item := range items {
		rows = append(rows, itemToRow(item, ctx.Config.ShowIcons, ctx.Current))
	}
	ctx.Table.SetRows(rows)
}

// itemToRow converts a single item to a table row.
func itemToRow(item vfs.FileSystemItem, showIcons bool, parent *vfs.Dir) table.Row {
	switch item.(type) {
	case *vfs.File, *vfs.Dir:
		var parentSize vfs.Size
		if parent != nil {
			parentSize = parent.Size()
		}
		return buildRow(item.Name(), item.Icon(), item.Size(), parentSize, showIcons)
	case *vfs.UpLink:
		return table.Row{"", "", "   ..  "}
	}
	return table.Row{}
}

// buildRow formats a single data row for the table.
func buildRow(name, emoji string, size, parentSize vfs.Size, showIcons bool) table.Row {
	percent := size.PercentOf(parentSize)
	color := gradientColor(percent)
	style := lipgloss.NewStyle().Foreground(color)
	if showIcons {
		name = emoji + " " + name
	}
	return table.Row{
		style.Render(fmt.Sprintf("%5.2f%%", percent)),
		style.Render(size.String() + " "),
		name + " ",
	}
}

// gradientColor returns a lipgloss color from green (0%) → yellow (50%) → red (100%).
func gradientColor(percent float64) color.Color {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	var r, g int
	if percent <= 50 {
		f := percent / 50.0
		r = int(math.Round(255 * f))
		g = 255
	} else {
		f := (percent - 50) / 50.0
		r = 255
		g = int(math.Round(255 * (1 - f)))
	}
	return lipgloss.Color(fmt.Sprintf("#%02x%02x00", r, g))
}
