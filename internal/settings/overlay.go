package settings

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

func (settings *Settings) Overlay(ctx *core.Context) string {
	rows := renderRows(&ctx.Config, settings.cursor)

	content := strings.Join([]string{
		styles.BoldAccent.Render("Settings"),
		"",
		rows,
	}, "\n")

	return styles.DialogBorder.Render(content)
}

func renderRows(cfg *config.Config, cursor int) string {
	rows := rowDefs(cfg)
	var b strings.Builder
	for i, r := range rows {
		marker := "  "
		nameStyle := styles.Muted
		if i == cursor {
			marker = styles.BoldAccent.Render("▶ ")
			nameStyle = styles.BoldAccent
		}
		name := nameStyle.Render(padRight(r.name, 12))
		value := styles.Main.Render(r.value)
		fmt.Fprintf(&b, "%s%s  %s", marker, name, value)
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func rowDefs(cfg *config.Config) []struct {
	name  string
	value string
} {
	return []struct {
		name  string
		value string
	}{
		{"sort", sortLabel(cfg.SortBy)},
		{"order", orderLabel(cfg.SortOrder)},
		{"group", config.GroupingString(cfg.Grouping)},
		{"icons", core.BoolLabel(cfg.ShowIcons)},
		{"hidden", core.BoolLabel(cfg.ShowHidden)},
		{"live filter", core.BoolLabel(cfg.LiveFilter)},
		{"notify level", cfg.NotifyLevel.PickLabel()},
		{"notify time", cfg.NotifyTimeout.PickLabel()},
	}
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func sortLabel(s config.SortBy) string {
	if s == config.SortByName {
		return "name"
	}
	return "size"
}

func orderLabel(o config.SortOrder) string {
	if o == config.Ascending {
		return "asc"
	}
	return "desc"
}
