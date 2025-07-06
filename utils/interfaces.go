package utils

import (
	"go-drive/ui"
	"sort"
	"strings"
	"unicode/utf8"

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
	style []lipgloss.Style
	pos   lipgloss.Position
}

func Stringize(f FileItem, path string, selected bool) string {
	om := orderedmap.New[string, column]()
	om.Set(
		"sync",
		column{2, []lipgloss.Style{
			lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
			lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
			lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		}, lipgloss.Left})
	om.Set("icon",
		column{2, []lipgloss.Style{
			lipgloss.NewStyle().Foreground(lipgloss.Color("170")),
			lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		}, lipgloss.Left},
	)
	om.Set("name", column{60, []lipgloss.Style{
		lipgloss.NewStyle(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Underline(true),
	}, lipgloss.Left})
	om.Set("size", column{10, []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}, lipgloss.Right})
	om.Set("shared", column{3, []lipgloss.Style{}, lipgloss.Left})

	dots_style := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	item_str := ""
	for pair := om.Oldest(); pair != nil; pair = pair.Next() {
		field := pair.Key
		col := pair.Value
		symbol := ""
		style := lipgloss.NewStyle()
		if len(col.style) > 0 {
			style = col.style[0]
		}

		switch field {

		case "icon":
			icon, ok := ui.MimeIcons[f.GetMimeType()]
			if !ok {
				icon = ui.MimeIcons["other"]
			}
			symbol = icon
			if selected {
				style = col.style[0]
			} else {
				style = col.style[1]
			}

		case "sync":
			if f.IsRemote() && f.IsLocal() {
				if f.IsDirectory() {
					symbol = " "
				} else {
					symbol = ui.StatusIcons["synced"]
					style = col.style[0]
				}
			} else if f.IsRemote() {
				symbol = ui.StatusIcons["remote"]
				style = col.style[1]
			} else if f.IsLocal() {
				symbol = ui.StatusIcons["local"]
				style = col.style[2]
			} else {
				panic("file is neither remote nor local")
			}

		case "name":
			name := f.GetName()
			if utf8.RuneCountInString(name) > col.width {
				name = name[:col.width-3] + "..."
			}
			dots := ""
			if !f.IsDirectory() {
				dots = dots_style.Render(strings.Repeat(".", col.width-utf8.RuneCountInString(name)))
			}
			if selected {
				name = col.style[1].Render(name)
			} else {
				name = col.style[0].Render(name)
			}
			symbol = name + dots

		case "shared":
			if !f.GetOwnedByMe() {
				symbol = ui.StatusIcons["shared"]
			}

		case "size":
			if !f.IsDirectory() {
				symbol = humanize.Bytes(f.GetSize())
			}
		}

		item_str = appendColumn(item_str, symbol, col.width, col.pos, style)
	}
	return item_str
}

func StringizeAll(f []FileItem, path string) []string {
	items := make([]string, len(f))
	for i, item := range f {
		items[i] = Stringize(item, path, false)
	}
	return items
}
