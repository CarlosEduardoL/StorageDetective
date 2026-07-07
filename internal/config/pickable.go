package config

import "fmt"

// Pickeable is a value that can be selected from a choose dialog. It knows how
// to display itself and what other values are available.
type Pickeable interface {
	PickLabel() string
	Options() []Pickeable
}

func (g Grouping) PickLabel() string { return GroupingString(g) }

func (Grouping) Options() []Pickeable {
	return []Pickeable{Mixed, FilesFirst, DirsFirst, FilesOnly, DirsOnly}
}

func (n NotifyLevel) PickLabel() string { return NotifyLevelString(n) }
func (NotifyLevel) Options() []Pickeable {
	return []Pickeable{NotifyAll, NotifyWarn, NotifyError, NotifyOff}
}

// NotifyTimeout is a toast timeout duration in seconds.
type NotifyTimeout int

func (t NotifyTimeout) PickLabel() string { return fmt.Sprintf("%ds", t) }

func (NotifyTimeout) Options() []Pickeable {
	return []Pickeable{NotifyTimeout(2), NotifyTimeout(3), NotifyTimeout(5), NotifyTimeout(10)}
}
