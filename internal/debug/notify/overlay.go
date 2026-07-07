package notify

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

func (m *Mode) Overlay(ctx *core.Context) string {
	var lines []string
	lines = append(lines, styles.BoldAccent.Render(" Debug: Notify "))
	lines = append(lines, "")

	fields := []struct {
		label string
		value string
	}{
		{"Title", m.title.View()},
		{"Detail", m.detail.View()},
		{"Severity", severityLabel(m.severity)},
	}

	for i, f := range fields {
		marker := "  "
		nameStyle := styles.Muted
		if i == m.focus {
			marker = styles.BoldAccent.Render("▸ ")
			nameStyle = styles.BoldAccent
		}
		lines = append(lines, fmt.Sprintf("%s%s  %s", marker, nameStyle.Render(f.label), f.value))
	}

	content := strings.Join(lines, "\n")
	width := min(45, ctx.Width-4)
	return styles.DialogBorder.Width(width).Render(content)
}

func severityLabel(s core.NotifySeverity) string {
	return s.String()
}
