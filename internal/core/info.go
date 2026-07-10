package core

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/layout"
	"github.com/SolracHQ/stex/internal/styles"
	"github.com/SolracHQ/stex/internal/vfs"

	"charm.land/lipgloss/v2"
)

// FileInfo holds the metadata that the right pane shows for a single file. The fields mirror
// the read only stat information plus the MIME type, which is detected from the first 512 bytes
// for regular files.
type FileInfo struct {
	Name          string
	Extension     string
	Size          vfs.Size
	ModTime       string
	Permissions   string
	IsSymlink     bool
	SymlinkTarget string
	MimeType      string
	IsDir         bool
}

// NewFileInfo reads the metadata of the file at path. Returns nil when the file does not exist
// or is not readable. The returned value is always safe to read, even for the nil case.
func NewFileInfo(path string) *FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	out := &FileInfo{
		Name:        filepath.Base(path),
		Extension:   filepath.Ext(path),
		Size:        vfs.Size(info.Size()),
		ModTime:     info.ModTime().Format("02 Jan 2006 15:04"),
		Permissions: info.Mode().String(),
		IsDir:       info.IsDir(),
	}

	if info.Mode()&os.ModeSymlink != 0 {
		out.IsSymlink = true
		if target, err := os.Readlink(path); err == nil {
			out.SymlinkTarget = target
		}
	}

	if !info.IsDir() && info.Mode().IsRegular() {
		out.MimeType = detectMIME(path)
	}

	return out
}

// detectMIME reads the first 512 bytes of path and returns the MIME type detected by net/http.
// Returns an empty string if the file cannot be opened or read.
func detectMIME(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	return http.DetectContentType(buf[:n])
}

// RenderFileInfo renders the right pane for a file.
func RenderFileInfo(info *FileInfo, bounds layout.Rect) string {
	if info == nil {
		return styles.CenterBox("No info available", bounds.Width, bounds.Height)
	}

	var lines []string
	lines = append(lines, styles.Bold(" File Info "))
	lines = append(lines, "")

	lines = append(lines, fieldLabel("Name:")+info.Name)
	if info.Extension != "" {
		lines = append(lines, fieldLabel("Extension:")+info.Extension)
	}
	if info.MimeType != "" {
		lines = append(lines, fieldLabel("Type:")+info.MimeType)
	}
	lines = append(lines, fieldLabel("Size:")+info.Size.String())
	lines = append(lines, fieldLabel("Modified:")+info.ModTime)
	lines = append(lines, fieldLabel("Permissions:")+info.Permissions)
	if info.IsSymlink {
		lines = append(lines, fieldLabel("Symlink:")+info.SymlinkTarget)
	}

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().Width(bounds.Width).Height(bounds.Height).Render(content)
}

// RenderDirInfo renders the right pane for a directory.
func RenderDirInfo(dir *vfs.Dir, info *FileInfo, bounds layout.Rect) string {
	if info == nil {
		return styles.CenterBox("No info available", bounds.Width, bounds.Height)
	}

	dirSize := dir.Size()
	count := dir.Count()
	children := TopChildren(dir, MaxTopChildren)

	var lines []string
	lines = append(lines, styles.Bold(" Directory Info "))
	lines = append(lines, "")

	lines = append(lines, fieldLabel("Name:")+info.Name)
	lines = append(lines, fieldLabel("Size:")+dirSize.String())
	lines = append(lines, fieldLabel("Files:")+fmt.Sprintf("%d", count.Files))
	lines = append(lines, fieldLabel("Dirs:")+fmt.Sprintf("%d", count.Dirs))
	lines = append(lines, fieldLabel("Modified:")+info.ModTime)
	lines = append(lines, fieldLabel("Permissions:")+info.Permissions)

	if len(children) > 0 {
		lines = append(lines, "")
		lines = append(lines, styles.Bold(" Largest Children "))
		lines = append(lines, "")
		for _, item := range children {
			pct := item.Size().PercentOf(dirSize)
			lines = append(lines, fmt.Sprintf(" %5.2f%% %10s  %s", pct, item.Size().String(), item.Name()))
		}
	}

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().Width(bounds.Width).Height(bounds.Height).Render(content)
}

// fieldLabel renders a right-pane field label in a fixed 14-cell column.
func fieldLabel(label string) string {
	return " " + lipgloss.NewStyle().Width(14).Align(lipgloss.Left).Render(label)
}

// MaxTopChildren is the number of children shown in the Largest Children section.
const MaxTopChildren = 10

// UpdateInfo caches the file stat for the current cursor position.
func UpdateInfo(ctx *Context) {
	if ctx.Current == nil {
		return
	}
	idx := ctx.Table.Cursor()
	if idx < 0 || idx >= len(ctx.Items) {
		return
	}
	item := ctx.Items[idx]
	path := item.FullPath()

	if _, ok := item.(*vfs.UpLink); ok {
		ctx.Info.file = nil
		return
	}

	if path != ctx.Info.Path {
		ctx.Info.Path = path
		ctx.Info.file = NewFileInfo(path)
	}
}

// InfoContent renders the info panel for the current cursor position. Returns empty when
// there is nothing to show (UpLink, nil file, or bounds are zero).
func InfoContent(ctx *Context) string {
	if ctx.Info.file == nil || ctx.Info.Bounds.Width == 0 || ctx.Info.Bounds.Height == 0 {
		return ""
	}
	idx := ctx.Table.Cursor()
	if idx < 0 || idx >= len(ctx.Items) {
		return ""
	}
	item := ctx.Items[idx]
	if _, ok := item.(*vfs.UpLink); ok {
		return ""
	}
	if dir, ok := item.(*vfs.Dir); ok {
		return RenderDirInfo(dir, ctx.Info.file, ctx.Info.Bounds)
	}
	return RenderFileInfo(ctx.Info.file, ctx.Info.Bounds)
}

// TopChildren returns the n largest direct children of dir, sorted by size descending.
func TopChildren(dir *vfs.Dir, n int) []vfs.FileSystemItem {
	if dir == nil || n <= 0 {
		return nil
	}
	cfg := config.Config{
		SortBy:    config.SortBySize,
		SortOrder: config.Descending,
		Grouping:  config.Mixed,
	}
	items := dir.ComputeItems(cfg)
	if len(items) > n {
		items = items[:n]
	}
	return items
}
