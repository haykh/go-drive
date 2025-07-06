package utils

import (
	"go-drive/ui"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	orderedmap "github.com/wk8/go-ordered-map/v2"
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

type column struct {
	width int
	pos   lipgloss.Position
}

func Stringize(f FileItem, path string) string {
	icon, ok := ui.MimeIcons[f.GetMimeType()]
	if !ok {
		icon = ui.MimeIcons["other"]
	}

	om := orderedmap.New[string, column]()
	om.Set("icon", column{2, lipgloss.Left})
	om.Set("sync", column{2, lipgloss.Left})
	om.Set("name", column{60, lipgloss.Left})
	om.Set("size", column{10, lipgloss.Right})
	om.Set("shared", column{3, lipgloss.Left})

	item_str := ""
	for pair := om.Oldest(); pair != nil; pair = pair.Next() {
		field := pair.Key
		col := pair.Value
		symbol := ""

		switch field {

		case "icon":
			symbol = icon

		case "sync":
			if f.IsRemote() && f.IsLocal() {
				if f.IsDirectory() {
					symbol = " "
				} else {
					symbol = ui.StatusIcons["synced"]
				}
			} else if f.IsRemote() {
				symbol = ui.StatusIcons["remote"]
			} else if f.IsLocal() {
				symbol = ui.StatusIcons["local"]
			} else {
				panic("file is neither remote nor local")
			}

		case "name":
			name := f.GetName()
			if len(name) > col.width {
				name = name[:col.width-3] + "..."
			}
			symbol = name

		case "shared":
			if !f.GetOwnedByMe() {
				symbol = ui.StatusIcons["shared"]
			}

		case "size":
			if !f.IsDirectory() {
				symbol = humanize.Bytes(f.GetSize())
			}
		}

		item_str = appendColumn(item_str, symbol, col.width, col.pos, lipgloss.NewStyle())
	}
	return item_str
}

func StringizeAll(f []FileItem, path string) []string {
	items := make([]string, len(f))
	for i, item := range f {
		items[i] = Stringize(item, path)
	}
	return items
}
