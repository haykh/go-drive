package browser

import (
	"fmt"
	"go-drive/utils"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type itemRenderer struct {
	filelist []utils.FileItem
	cwd      []string
}

func (d itemRenderer) CWD() string {
	return strings.Join(d.cwd, "/")
}

func (d itemRenderer) Height() int                             { return 1 }
func (d itemRenderer) Spacing() int                            { return 0 }
func (d itemRenderer) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemRenderer) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	str := utils.Stringize(d.filelist[index], d.CWD())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, "> ", strings.Join(s, " ")))
		}
	}

	fmt.Fprint(w, fn(str))
}
