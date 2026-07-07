package choose

import (
	"strings"

	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

func (ch *Choose) Overlay(ctx *core.Context) string {
	width := max(1, min(60, ctx.Width-4))

	var rows []string
	for i, opt := range ch.options {
		if i == ch.cursor {
			rows = append(rows, styles.BoldAccent.Render("▶ "+opt.Label))
		} else {
			rows = append(rows, styles.Main.Render("  "+opt.Label))
		}
	}

	content := strings.Join([]string{
		styles.BoldAccent.Render(ch.title),
		"",
		strings.Join(rows, "\n"),
	}, "\n")

	return styles.DialogBorder.Width(width).Render(content)
}
