package utils

import (
	"go-drive/ui"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

type FileManager interface {
	GetFileList(string, bool) ([]FileItem, error)
}

type FileItem interface {
	IsDirectory() bool
	IsPDF() bool
	IsUnrecognized() bool

	IsLocal() bool
	IsRemote() bool

	GetName() string
	GetSize() uint64
	GetMimeType() string
	GetModifiedTime() string
	GetOwnedByMe() bool
}

func Sorted(f []FileItem) []FileItem {
	sort.Slice(f, func(i, j int) bool {
		if f[i].GetMimeType() == f[j].GetMimeType() {
			return f[i].GetName() < f[j].GetName()
		} else if f[i].IsDirectory() || f[j].IsDirectory() {
			return f[i].IsDirectory()
		} else if f[i].IsPDF() || f[j].IsPDF() {
			return f[i].IsPDF()
		} else if f[i].GetMimeType() == "other" || f[j].GetMimeType() == "other" {
			return f[j].GetMimeType() == "other"
		}
		return f[i].GetMimeType() < f[j].GetMimeType()
	})
	return f
}

func appendColumn(current, column string, width int, pos lipgloss.Position, style lipgloss.Style) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		current,
		style.Width(width).Align(pos).Render(column),
	)
}

func Stringize(f FileItem, path string, cols []string) string {
	icon, ok := ui.MimeIcons[f.GetMimeType()]
	if !ok {
		icon = ui.MimeIcons["other"]
	}
	if len(cols) == 0 {
		cols = []string{"icon", "sync", "name", "size", "shared"}
	}

	style := lipgloss.NewStyle()
	item_str := ""
	for _, field := range cols {
		switch field {
		case "icon":
			item_str = appendColumn(item_str, icon, 2, lipgloss.Left, style)
		case "sync":
			if f.IsRemote() && !f.IsDirectory() {
				item_str = appendColumn(item_str, "", 2, lipgloss.Left, style)
			} else {
				item_str = appendColumn(item_str, "", 2, lipgloss.Left, style)
			}
		case "name":
			name := f.GetName()
			if len(name) > 60 {
				name = name[:57] + "..."
			}
			item_str = appendColumn(item_str, name, 60, lipgloss.Left, style)
		case "shared":
			if !f.GetOwnedByMe() {
				item_str = appendColumn(item_str, " ", 2, lipgloss.Left, style)
			} else {
				item_str = appendColumn(item_str, "", 2, lipgloss.Left, style)
			}
		case "size":
			if !f.IsDirectory() {
				item_str = appendColumn(item_str, humanize.Bytes(f.GetSize()), 10, lipgloss.Right, style)
			}
		}
	}
	return item_str
}

func StringizeAll(f []FileItem, path string, format []string) []string {
	items := make([]string, len(f))
	for i, item := range f {
		items[i] = Stringize(item, path, format)
	}
	return items
}
