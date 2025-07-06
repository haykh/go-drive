package browser

import (
	"bytes"
	"go-drive/filesystem"
	"go-drive/ui"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func FileBrowser(
	mgr filesystem.FileManager,
	path string,
	debug_mode bool,
	debugBuffer *bytes.Buffer,
) error {
	l := list.New([]list.Item{}, itemRenderer{nil, []int{}, []string{}}, 20, 30)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = ui.BrowserPagination
	model := browserModel{
		debug_mode: debug_mode,
		format:     []string{},

		cwd:         strings.Split(path, "/"),
		filemanager: mgr,
		filelist:    []filesystem.FileItem{},
		filesinsync: []int{},

		help:    help.New(),
		list:    l,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),

		loading:     true,
		status:      "",
		quitting:    false,
		debugLines:  []string{},
		debugBuffer: debugBuffer,

		keys: keys,
	}

	p := tea.NewProgram(model)

	_, err := p.Run()
	return err
}
