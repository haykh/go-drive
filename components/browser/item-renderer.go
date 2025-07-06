package browser

import (
	"fmt"
	"go-drive/filesystem"
	"go-drive/ui"
	"io"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type itemRenderer struct {
	filelist []filesystem.FileItem
	syncing  []int
	cwd      []string
}

func (d itemRenderer) CWD() string {
	return strings.Join(d.cwd, "/")
}

func (d itemRenderer) Height() int                             { return 1 }
func (d itemRenderer) Spacing() int                            { return 0 }
func (d itemRenderer) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemRenderer) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	is_syncing := slices.Contains(d.syncing, index)
	is_hovered := m.Index() == index
	str := filesystem.Stringize(d.filelist[index], d.CWD(), is_hovered, is_syncing)

	fn := ui.BrowserItem.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return ui.BrowserSelectedItem.Render(
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					"> ",
					strings.Join(s, " "),
				))
		}
	}

	fmt.Fprint(w, fn(str))
}

func FileListToItems(r []filesystem.FileItem) []list.Item {
	items := []list.Item{}
	for i := range r {
		items = append(items, browserItem(fmt.Sprintf("%d", i)))
	}
	return items
}

type browserItem string

func (i browserItem) FilterValue() string { return "" }
