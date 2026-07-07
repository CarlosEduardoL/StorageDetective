package app

import (
	"strings"
	"time"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
	"github.com/charmbracelet/x/ansi"

	tea "charm.land/bubbletea/v2"
)

const (
	notifyBoxWidth  = 40
	notifyBoxHeight = 6
	notifyMargin    = 2
)

// RenderNotify renders the notification box with the given text and severity.
func RenderNotify(text, detail string, severity core.NotifySeverity) string {
	severityLabel := severity.String()

	topBorder := renderTopBorder(severityLabel, severity)
	bottomBorder := "╰" + strings.Repeat("─", notifyBoxWidth-2) + "╯"

	innerWidth := notifyBoxWidth - 4

	textLine := wrapCenter(text, innerWidth, notifyBoxWidth)
	detailLine := ""
	if detail != "" {
		detailLine = wrapCenter(styles.Dim.Render(detail), innerWidth, notifyBoxWidth)
	}

	contentLines := make([]string, 0, notifyBoxHeight)
	contentLines = append(contentLines, topBorder)
	contentLines = append(contentLines, "│"+strings.Repeat(" ", notifyBoxWidth-2)+"│")
	contentLines = append(contentLines, "│"+textLine+"│")
	if detailLine != "" {
		contentLines = append(contentLines, "│"+detailLine+"│")
	} else {
		contentLines = append(contentLines, "│"+strings.Repeat(" ", notifyBoxWidth-2)+"│")
	}
	contentLines = append(contentLines, "│"+strings.Repeat(" ", notifyBoxWidth-2)+"│")
	contentLines = append(contentLines, bottomBorder)

	return strings.Join(contentLines, "\n")
}

// OverlayNotify composites the notification box over the background content, placing it in a
// fixed position so it is visible without displacing the base view.
func OverlayNotify(background string, toast string) string {
	if toast == "" || background == "" {
		return background
	}

	bgLines := strings.Split(background, "\n")
	bgWidth := 0
	for _, line := range bgLines {
		if w := ansi.StringWidth(line); w > bgWidth {
			bgWidth = w
		}
	}
	if bgWidth == 0 {
		return background
	}

	toastLines := strings.Split(toast, "\n")
	toastWidth := 0
	for _, line := range toastLines {
		if w := ansi.StringWidth(line); w > toastWidth {
			toastWidth = w
		}
	}
	toastHeight := len(toastLines)

	offsetX := max(bgWidth-toastWidth-notifyMargin, 0)
	offsetY := notifyMargin

	var buf strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			buf.WriteByte('\n')
		}
		if i < offsetY || i >= offsetY+toastHeight {
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

		toastLine := toastLines[i-offsetY]
		buf.WriteString(toastLine)
		pos += ansi.StringWidth(toastLine)

		right := ansi.TruncateLeft(bgLine, pos, "")
		bgW := ansi.StringWidth(bgLine)
		rightW := ansi.StringWidth(right)
		if rightW <= bgW-pos {
			buf.WriteString(strings.Repeat(" ", bgW-rightW-pos))
		}
		buf.WriteString(right)
	}

	if len(bgLines) < offsetY+toastHeight {
		for i := len(bgLines); i < offsetY+toastHeight; i++ {
			buf.WriteByte('\n')
			left := strings.Repeat(" ", offsetX)
			toastLine := toastLines[i-len(bgLines)]
			buf.WriteString(left)
			buf.WriteString(toastLine)
		}
	}

	return buf.String()
}

func renderTopBorder(severity string, sev core.NotifySeverity) string {
	left := "╭── " + severityColor(severity, sev) + " "
	right := " [× esc] "
	sep := "─"
	fill := max(notifyBoxWidth-ansi.StringWidth(left)-ansi.StringWidth(right)-2, 1)
	sepLine := strings.Repeat(sep, fill)
	right = sepLine + right + "─╮"
	return left + right
}

func severityColor(label string, sev core.NotifySeverity) string {
	switch sev {
	case core.NotifyInfo:
		return styles.NotifyInfo.Render(label)
	case core.NotifyWarn:
		return styles.NotifyWarn.Render(label)
	case core.NotifyError:
		return styles.NotifyError.Render(label)
	default:
		return styles.NotifyInfo.Render(label)
	}
}

func wrapCenter(text string, innerWidth, boxWidth int) string {
	text = ansi.Truncate(text, innerWidth, "…")
	padding := max(boxWidth-2-ansi.StringWidth(text), 0)
	leftPad := padding / 2
	rightPad := padding - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// enqueueNotify adds a notification to the end of the queue and returns a tick command for auto
// dismissal. If the queue already has a notification showing the tick is handled by
// dismissNotify.
func (app *App) enqueueNotify(msg core.AppNotify) tea.Cmd {
	app.notifyQueue = append(app.notifyQueue, msg)
	if len(app.notifyQueue) == 1 {
		return core.Tick(notifyTimeout(app.ctx.Config), core.AppNotifyDone{})
	}
	return nil
}

// dismissNotify removes the current notification and returns a tick for the next queued
// notification, or nil when the queue is empty.
func (app *App) dismissNotify() tea.Cmd {
	if len(app.notifyQueue) == 0 {
		return nil
	}
	app.notifyQueue = app.notifyQueue[1:]
	if len(app.notifyQueue) > 0 {
		return core.Tick(notifyTimeout(app.ctx.Config), core.AppNotifyDone{})
	}
	return nil
}

func notifyTimeout(cfg config.Config) time.Duration {
	return time.Duration(int(cfg.NotifyTimeout)) * time.Second
}

// notifyAllowed checks whether a notification of the given severity should be shown based on the
// user's NotifyLevel setting.
func (app *App) notifyAllowed(severity core.NotifySeverity) bool {
	switch app.ctx.Config.NotifyLevel {
	case config.NotifyOff:
		return false
	case config.NotifyError:
		return severity == core.NotifyError
	case config.NotifyWarn:
		return severity == core.NotifyError || severity == core.NotifyWarn
	default:
		return true
	}
}
