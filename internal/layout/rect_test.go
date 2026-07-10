package layout

import "testing"

func TestContains(t *testing.T) {
	r := New(5, 10, 100, 50)
	cases := []struct {
		x, y int
		want bool
	}{
		{5, 10, true},
		{104, 59, true},
		{4, 10, false},
		{5, 9, false},
		{105, 10, false},
		{5, 60, false},
	}
	for _, c := range cases {
		got := r.Contains(c.x, c.y)
		if got != c.want {
			t.Fatalf("Contains(%d,%d) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

func TestShrink(t *testing.T) {
	r := New(0, 0, 80, 24).Shrink(1)
	if r.Width != 78 || r.Height != 22 {
		t.Fatalf("Shrink(1): want 78x22, got %dx%d", r.Width, r.Height)
	}
}

func TestSplitVSumsToTotal(t *testing.T) {
	r := New(0, 0, 100, 20)
	l, r2 := r.SplitV(50)
	if l.Width+r2.Width != r.Width {
		t.Fatalf("SplitV 50%%: %d + %d = %d, want %d", l.Width, r2.Width, l.Width+r2.Width, r.Width)
	}
}

func TestTopAndBottom(t *testing.T) {
	r := New(0, 0, 80, 24)
	top, rest := r.Top(1)
	if top.Height != 1 || rest.Height != 23 {
		t.Fatalf("Top(1): top=%d rest=%d, want 1,23", top.Height, rest.Height)
	}
	_, bot := r.Bottom(1)
	if bot.Height != 1 {
		t.Fatalf("Bottom(1): bot=%d, want 1", bot.Height)
	}
	if bot.Y != 23 {
		t.Fatalf("Bottom(1): y=%d, want 23", bot.Y)
	}
}
