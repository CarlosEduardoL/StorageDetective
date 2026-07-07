package core

import tea "charm.land/bubbletea/v2"

// NotifySeverity indicates the notification level.
type NotifySeverity int

const (
	NotifyInfo NotifySeverity = iota
	NotifyWarn
	NotifyError
)

func (s NotifySeverity) String() string {
	switch s {
	case NotifyInfo:
		return "Info"
	case NotifyWarn:
		return "Warning"
	case NotifyError:
		return "Error"
	}
	return ""
}

// AppNotify requests the app to show a notification.
type AppNotify struct {
	Text     string
	Detail   string
	Severity NotifySeverity
}

// AppNotifyDone clears the current notification and shows the next in queue.
type AppNotifyDone struct{}

// NewNotifyCmd returns a tea.Cmd that emits an AppNotify for the given text and severity.
func NewNotifyCmd(text, detail string, severity NotifySeverity) tea.Cmd {
	return func() tea.Msg {
		return AppNotify{Text: text, Detail: detail, Severity: severity}
	}
}
