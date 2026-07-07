package debug

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

type menuAction struct {
	label string
	desc  string
}

var menuActions = []menuAction{
	{label: "Notify", desc: "Fire a notification with custom text and severity"},
}

func (m *Mode) Overlay(ctx *core.Context) string {
	var lines []string
	lines = append(lines, styles.BoldAccent.Render(" Debug Mode "))
	lines = append(lines, "")

	for i, action := range menuActions {
		marker := "  "
		nameStyle := styles.Muted
		if i == m.cursor {
			marker = styles.BoldAccent.Render("▶ ")
			nameStyle = styles.BoldAccent
		}
		lines = append(lines, fmt.Sprintf("%s%s", marker, nameStyle.Render(action.label)))
		lines = append(lines, fmt.Sprintf("   %s", styles.Dim.Render(action.desc)))
	}

	content := strings.Join(lines, "\n")
	return styles.DialogBorder.Render(content)
}
