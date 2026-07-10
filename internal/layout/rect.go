// Package layout provides integer rectangle operations for screen subdivision and layout
// arithmetic. Every value is a plain struct with exported fields, safe to copy.
package layout

// Rect is a 2D axis-aligned rectangle with integer coordinates and dimensions.
type Rect struct {
	X, Y, Width, Height int
}

// New returns a Rect at (x,y) with the given dimensions.
func New(x, y, w, h int) Rect {
	return Rect{X: x, Y: y, Width: w, Height: h}
}

// Contains reports whether (x,y) is inside the rectangle, inclusive of the
// bottom and right edges.
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x <= r.X+r.Width-1 &&
		y >= r.Y && y <= r.Y+r.Height-1
}

// Shrink returns a rect shrunk by n on all four sides. A negative n expands.
func (r Rect) Shrink(n int) Rect {
	return Rect{
		X:      r.X + n,
		Y:      r.Y + n,
		Width:  max(r.Width-2*n, 0),
		Height: max(r.Height-2*n, 0),
	}
}

// ShrinkSides returns a rect shrunk by different amounts per side.
func (r Rect) ShrinkSides(left, right, top, bottom int) Rect {
	return Rect{
		X:      r.X + left,
		Y:      r.Y + top,
		Width:  max(r.Width-left-right, 0),
		Height: max(r.Height-top-bottom, 0),
	}
}

// Top carves n lines from the top of r and returns the carved rect and the
// remainder. n must be non-negative and no larger than r.Height.
func (r Rect) Top(n int) (Rect, Rect) {
	if n > r.Height {
		n = r.Height
	}
	return r.ShrinkSides(0, 0, 0, r.Height-n),
		Rect{X: r.X, Y: r.Y + n, Width: r.Width, Height: r.Height - n}
}

// Bottom carves n lines from the bottom of r and returns the remainder and
// the carved rect. n must be non-negative and no larger than r.Height.
func (r Rect) Bottom(n int) (Rect, Rect) {
	if n > r.Height {
		n = r.Height
	}
	return Rect{X: r.X, Y: r.Y, Width: r.Width, Height: r.Height - n},
		Rect{X: r.X, Y: r.Y + r.Height - n, Width: r.Width, Height: n}
}

// SplitV splits r vertically at pct percent of the width.
func (r Rect) SplitV(pct int) (Rect, Rect) {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	leftW := r.Width * pct / 100
	rightW := r.Width - leftW
	return Rect{X: r.X, Y: r.Y, Width: leftW, Height: r.Height},
		Rect{X: r.X + leftW, Y: r.Y, Width: rightW, Height: r.Height}
}

// SplitH splits r horizontally at pct percent of the height.
func (r Rect) SplitH(pct int) (Rect, Rect) {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	topH := r.Height * pct / 100
	botH := r.Height - topH
	return Rect{X: r.X, Y: r.Y, Width: r.Width, Height: topH},
		Rect{X: r.X, Y: r.Y + topH, Width: r.Width, Height: botH}
}
