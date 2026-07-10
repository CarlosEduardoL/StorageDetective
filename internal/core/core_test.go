package core

import (
	"testing"
)

func TestRenderBaseEmptyContext(t *testing.T) {
	ctx := &Context{}
	if got := RenderBase(ctx); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
