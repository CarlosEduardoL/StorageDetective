package settings

import (
	"strings"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"

	"charm.land/lipgloss/v2"
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
	var rendered []string
	for i, r := range rows {
		marker := "  "
		nameStyle := styles.Muted
		if i == cursor {
			marker = styles.BoldAccent.Render("▶ ")
			nameStyle = styles.BoldAccent
		}
		name := nameStyle.Render(r.name)
		padded := lipgloss.NewStyle().Width(12).Render(name)
		value := styles.Main.Render(r.value)
		rendered = append(rendered, marker+padded+"  "+value)
	}
	return strings.Join(rendered, "\n")
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
		{"glyphs", core.BoolLabel(cfg.ShowPowerGlyphs)},
		{"hidden", core.BoolLabel(cfg.ShowHidden)},
		{"live filter", core.BoolLabel(cfg.LiveFilter)},
		{"notify level", cfg.NotifyLevel.PickLabel()},
		{"notify time", cfg.NotifyTimeout.PickLabel()},
	}
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
