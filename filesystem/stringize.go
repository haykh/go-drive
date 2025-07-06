package filesystem

import (
	"go-drive/ui"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

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

func Stringize(f FileItem, path string, selected, syncing bool) string {
	om := orderedmap.New[string, column]()
	om.Set("sync", column{2, lipgloss.Left})
	om.Set("icon", column{2, lipgloss.Left})
	om.Set("name", column{60, lipgloss.Left})
	om.Set("size", column{10, lipgloss.Right})
	om.Set("shared", column{3, lipgloss.Right})

	dots_style := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	item_str := ""
	for pair := om.Oldest(); pair != nil; pair = pair.Next() {
		field := pair.Key
		col := pair.Value
		symbol := ""
		style := lipgloss.NewStyle()

		switch field {

		case "icon":
			icon, ok := ui.MimeIcons[f.GetMimeType()]
			if !ok {
				icon = ui.MimeIcons["other"]
			}
			symbol = icon
			if selected {
				style = ui.MimeIconSelectedStyle
			} else {
				style = ui.MimeIconStyle
			}

		case "sync":
			if syncing {
				symbol = ui.StatusIcons["syncing"]
				style = ui.SyncingStyle
			} else if f.IsRemote() && f.IsLocal() {
				if f.IsDirectory() {
					symbol = " "
				} else {
					symbol = ui.StatusIcons["synced"]
					style = ui.SyncedStyle
				}
			} else if f.IsRemote() {
				symbol = ui.StatusIcons["remote"]
				style = ui.RemoteStyle
			} else if f.IsLocal() {
				symbol = ui.StatusIcons["local"]
				style = ui.LocalStyle
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
				name = ui.FileNameSelectedStyle.Render(name)
			} else {
				name = ui.FileNameStyle.Render(name)
			}
			symbol = name + dots

		case "shared":
			if !f.GetOwnedByMe() {
				symbol = ui.StatusIcons["shared"]
			}

		case "size":
			if !f.IsDirectory() {
				symbol = humanize.Bytes(f.GetSize())
				style = ui.FileSizeStyle
			}
		}

		item_str = appendColumn(item_str, symbol, col.width, col.pos, style)
	}
	return item_str
}

func StringizeAll(f []FileItem, path string) []string {
	items := make([]string, len(f))
	for i, item := range f {
		items[i] = Stringize(item, path, false, false)
	}
	return items
}
