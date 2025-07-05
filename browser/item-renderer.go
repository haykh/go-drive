package browser

import (
	"fmt"
	"go-drive/ui"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/drive/v3"
)

type itemRenderer struct {
	filelist *drive.FileList
	format   string
}

func (d itemRenderer) Height() int                             { return 1 }
func (d itemRenderer) Spacing() int                            { return 0 }
func (d itemRenderer) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemRenderer) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	file := d.filelist.Files[index]
	str := ui.StringizeItem(file, d.format)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
