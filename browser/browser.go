package browser

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/drive/v3"
)

func FileBrowser(srv *drive.Service, path string, debug_mode bool) error {
	l := list.New([]list.Item{}, itemRenderer{nil, ""}, 20, 40)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	model := browserModel{
		srv:        srv,
		debug_mode: debug_mode,
		format:     "%icon% %shared% %name% %size%",

		cwd: strings.Split(path, "/"),

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
