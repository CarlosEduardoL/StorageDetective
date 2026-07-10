// Package styles centralizes the visual constants and reusable style builders for every
// dialog, input, and display element in the application.
package styles

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

// DialogBorder is a double border style with the accent color and inner horizontal padding.
// The caller should set Width before rendering.
var DialogBorder = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color(AccentColor)).
	Padding(0, 2)

// DialogBorderWide is like DialogBorder but with wider horizontal padding.
var DialogBorderWide = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color(AccentColor)).
	Padding(0, 4)

// BorderNorm is a normal (single line) border with dim foreground.
var BorderNorm = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(DimColor))

// Dim is a dim foreground style for separators and less prominent text.
var Dim = lipgloss.NewStyle().Foreground(lipgloss.Color(DimColor))

// Accent renders text in accent color (yellow).
var Accent = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor))

// Active renders text in active color (cyan).
var Active = lipgloss.NewStyle().Foreground(lipgloss.Color(ActiveColor))

// BoldAccent renders text in accent color and bold.
var BoldAccent = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor)).Bold(true)

// Main renders text in main (white) color.
var Main = lipgloss.NewStyle().Foreground(lipgloss.Color(MainColor))

// BoldMain renders text in main color and bold.
var BoldMain = lipgloss.NewStyle().Foreground(lipgloss.Color(MainColor)).Bold(true)

// Muted renders text in muted (darker gray) color.
var Muted = lipgloss.NewStyle().Foreground(lipgloss.Color(MutedColor))

// NotifyInfo renders text in info (pale blue) color for info toasts.
var NotifyInfo = lipgloss.NewStyle().Foreground(lipgloss.Color(InfoColor))

// NotifyWarn renders text in accent (yellow) color for warning toasts.
var NotifyWarn = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor))

// NotifyError renders text in error (bright red) color for error toasts.
var NotifyError = lipgloss.NewStyle().Foreground(lipgloss.Color(ErrorColor))

// PowerbarMode is the style for the mode name segment in the power bar. Yellow background,
// black text for contrast.
var PowerbarMode = lipgloss.NewStyle().
	Background(lipgloss.Color(AccentColor)).
	Foreground(lipgloss.Color("0"))

// PowerbarSegment is the style for neutral segments in the power bar (grouping, size, counts).
// Dim background, main text.
var PowerbarSegment = lipgloss.NewStyle().
	Background(lipgloss.Color(DimColor)).
	Foreground(lipgloss.Color(MainColor))

// PowerbarSegmentAlt is the alternating style for odd neutral segments (hidden, counts).
// Slightly darker background than PowerbarSegment for a subtle zebra stripe effect.
var PowerbarSegmentAlt = lipgloss.NewStyle().
	Background(lipgloss.Color(DimAltColor)).
	Foreground(lipgloss.Color(MainColor))

// PowerbarPath is the style for the path segment in the power bar. Dim background, muted text
// so the path does not compete with the mode label.
var PowerbarPath = lipgloss.NewStyle().
	Background(lipgloss.Color(DimColor)).
	Foreground(lipgloss.Color(MutedColor))

// PowerbarSep is the style for segment separators in the power bar. The foreground is set to
// the same color as the left segment's background so the separator blends into it. The
// background is set to the next segment's background so it points into it. The zero value is
// safe to use when the next background is DimColor.
var PowerbarSep = lipgloss.NewStyle().
	Foreground(lipgloss.Color(AccentColor)).
	Background(lipgloss.Color(DimColor))

// PowerbarGlyph is the style for the powerline glyph between segments when ShowPowerGlyphs is
// true. The foreground renders in the left segment's background color and the background is the
// right segment's background, so the glyph blends into the left and points into the right.
var PowerbarGlyph = lipgloss.NewStyle().
	Foreground(lipgloss.Color(AccentColor)).
	Background(lipgloss.Color(DimColor))

// PowerbarSepToAlt transitions from a regular dim segment to an alt segment.
var PowerbarSepToAlt = lipgloss.NewStyle().
	Foreground(lipgloss.Color(DimColor)).
	Background(lipgloss.Color(DimAltColor))

// PowerbarSepFromAlt transitions from an alt segment back to a regular dim segment.
var PowerbarSepFromAlt = lipgloss.NewStyle().
	Foreground(lipgloss.Color(DimAltColor)).
	Background(lipgloss.Color(DimColor))

// PowerbarGlyphToAlt is the nerdfont version of PowerbarSepToAlt.
var PowerbarGlyphToAlt = lipgloss.NewStyle().
	Foreground(lipgloss.Color(DimColor)).
	Background(lipgloss.Color(DimAltColor))

// PowerbarGlyphFromAlt is the nerdfont version of PowerbarSepFromAlt.
var PowerbarGlyphFromAlt = lipgloss.NewStyle().
	Foreground(lipgloss.Color(DimAltColor)).
	Background(lipgloss.Color(DimColor))

// Bold wraps text in bold weight without changing color.
func Bold(text string) string {
	return lipgloss.NewStyle().Bold(true).Render(text)
}

// CenterBox renders text centred both horizontally and vertically within the given dimensions.
func CenterBox(text string, width, height int) string {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(text)
}

// CenterText formats a single line of text centred within a fixed width. Used for the help
// footer and any other one-line centred display.
func CenterText(text string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(text)
}

// SepString renders a separator glyph between two background colors. fromBg is the background
// of the segment the visual comes from, toBg is the background of the segment it points into.
// Either can be empty to leave that side transparent (the terminal background shows through).
func SepString(glyph, fromBg, toBg string) string {
	st := lipgloss.NewStyle()
	if fromBg != "" {
		st = st.Foreground(lipgloss.Color(fromBg))
	}
	if toBg != "" {
		st = st.Background(lipgloss.Color(toBg))
	}
	return st.Render(glyph)
}

// HelpKey is the style for the key part of the help footer.
var HelpKey = lipgloss.NewStyle().Foreground(lipgloss.Color(ActiveColor))

// HelpDesc is the style for the description part of the help footer.
var HelpDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(MainColor))

// HelpSep is the style for separators in the help footer.
var HelpSep = lipgloss.NewStyle().Foreground(lipgloss.Color(DimColor))

// HelpDefaults returns a help.Model with styles configured to match the application palette.
func HelpDefaults() help.Model {
	m := help.New()
	m.Styles.FullKey = HelpKey
	m.Styles.FullDesc = HelpDesc
	m.Styles.FullSeparator = HelpSep
	m.Styles.ShortKey = HelpKey
	m.Styles.ShortDesc = HelpDesc
	m.Styles.ShortSeparator = HelpSep
	return m
}

// TableDefault returns a table.Styles with a bold header and the project's SelectBg highlight.
func TableDefault() table.Styles {
	s := table.DefaultStyles()
	s.Header = lipgloss.NewStyle().Bold(true)
	s.Cell = lipgloss.NewStyle()
	s.Selected = lipgloss.NewStyle().Background(lipgloss.Color(SelectBg))
	return s
}
