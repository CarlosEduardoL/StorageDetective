package app

import (
	"testing"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/explorer"
	"github.com/SolracHQ/stex/internal/layout"
	"github.com/SolracHQ/stex/internal/styles"
	"github.com/SolracHQ/stex/internal/vfs"

	tea "charm.land/bubbletea/v2"
)

func TestAppNewBuildsContext(t *testing.T) {
	root := &vfs.Dir{}
	m := New(".", config.DefaultConfig(), root)
	a, ok := m.(*App)
	if !ok {
		t.Fatalf("expected *App, got %T", m)
	}
	if a.ctx.Root != root {
		t.Fatal("expected ctx.Root to match the passed root")
	}
	if a.ctx.Current != root {
		t.Fatal("expected ctx.Current to match the passed root")
	}
	if a.ctx.Path != "." {
		t.Fatalf("expected ctx.Path to be '.', got %q", a.ctx.Path)
	}
}

func TestAppFirstModeIsExplorer(t *testing.T) {
	root := &vfs.Dir{}
	m := New(".", config.DefaultConfig(), root)
	a := m.(*App)
	if a.mode == nil {
		t.Fatal("expected a mode to be set")
	}
	if _, ok := a.mode.(*explorer.Explorer); !ok {
		t.Fatalf("expected *explorer.Explorer, got %T", a.mode)
	}
}

func TestAppTableFocusedAtStart(t *testing.T) {
	root := &vfs.Dir{}
	m := New(".", config.DefaultConfig(), root)
	a := m.(*App)
	if !a.ctx.Table.Focused() {
		t.Fatal("expected table to be focused at start so key navigation works")
	}
}

func TestAppHandlesResizeBeforeInit(t *testing.T) {
	ctx := &core.Context{}
	ctx.SetScreen(layout.New(0, 0, 80, 24))
	a := &App{
		ctx:    ctx,
		screen: layout.New(0, 0, 80, 24),
		help:   styles.HelpDefaults(),
		mode:   &fakeMode{},
	}
	_, _ = a.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if a.ctx.Screen().Width != 118 {
		t.Fatalf("expected content width 118, got %d", a.ctx.Screen().Width)
	}

}

// emptyKeyMap is a help.KeyMap implementation with no bindings.
type emptyKeyMap struct{}

func (emptyKeyMap) ShortHelp() []key.Binding  { return []key.Binding{} }
func (emptyKeyMap) FullHelp() [][]key.Binding { return [][]key.Binding{} }

// fakeMode is a test mode that returns next from Update, used to verify the app installs the
// returned mode and runs its Init.
type fakeMode struct {
	next core.Mode
}

func (f *fakeMode) Init(_ *core.Context) tea.Cmd { return nil }

func (f *fakeMode) Update(_ *core.Context, _ tea.Msg) (core.Mode, tea.Cmd) {
	return f.next, nil
}

func (f *fakeMode) Overlay(_ *core.Context) string { return "" }

func (f *fakeMode) Help() help.KeyMap { return emptyKeyMap{} }
func (f *fakeMode) Name() string      { return "fake" }

func TestAppSwapsModeWhenReturned(t *testing.T) {
	ctx := &core.Context{}
	ctx.SetScreen(layout.New(0, 0, 80, 24))
	a := &App{
		ctx:    ctx,
		screen: layout.New(0, 0, 80, 24),
		help:   styles.HelpDefaults(),
		mode:   &fakeMode{next: &fakeMode{}},
	}
	_, _ = a.Update(nil)
	if _, ok := a.mode.(*fakeMode); !ok {
		t.Fatalf("expected mode swap to *fakeMode, got %T", a.mode)
	}
}

func TestOverlayCenterEmptyForeground(t *testing.T) {
	bg := "abc\ndef"
	if got := overlayCenter(bg, ""); got != bg {
		t.Fatalf("expected background returned, got %q", got)
	}
}

func TestOverlayCenterEmptyBackground(t *testing.T) {
	if got := overlayCenter("", "abc"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
