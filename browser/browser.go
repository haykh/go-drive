package browser

import (
	"go-drive/utils"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func FileBrowser(mgr utils.FileManager, path string, debug_mode bool) error {
	l := list.New([]list.Item{}, itemRenderer{nil, []string{}, ""}, 20, 40)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	model := browserModel{
		debug_mode: debug_mode,
		format:     "%icon% %sync% %shared% %name% %size%",

		cwd:         strings.Split(path, "/"),
		filemanager: mgr,

		help:    help.New(),
		list:    l,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),

		loading:  true,
		quitting: false,

		keys: keys,
	}

	p := tea.NewProgram(model)

	_, err := p.Run()
	return err
}
