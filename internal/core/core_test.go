package core

import (
	"strings"
	"testing"
)

func TestBlankEmptyDims(t *testing.T) {
	if Blank(0, 0) != "" {
		t.Fatal("expected empty for zero dims")
	}
	if Blank(10, 0) != "" {
		t.Fatal("expected empty for zero height")
	}
	if Blank(0, 5) != "" {
		t.Fatal("expected empty for zero width")
	}
}

func TestBlankRendersFullSize(t *testing.T) {
	out := Blank(4, 3)
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if len(line) != 4 {
			t.Fatalf("line %d: expected 4 chars, got %d", i, len(line))
		}
	}
}

func TestRenderBaseEmptyContext(t *testing.T) {
	ctx := &Context{}
	if got := RenderBase(ctx, nil); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
