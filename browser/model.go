package browser

import (
	"fmt"
	"go-drive/api"
	"go-drive/ui"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/api/drive/v3"
)

type browserItem string

func (i browserItem) FilterValue() string { return "" }

type doneLoadingMsg struct {
	filelist *drive.FileList
	err      error
}

type browserModel struct {
	// model configurations
	srv        *drive.Service
	debug_mode bool
	format     string

	// model state
	cwd      []string
	filelist *drive.FileList

	// ui components
	help    help.Model
	list    list.Model
	spinner spinner.Model

	// ui state
	loading  bool
	quitting bool

	// ui configurations
	keys keyMap

	// choice    string
	// choiceIdx int
}

func (m browserModel) CWD() string {
	return strings.Join(m.cwd, "/")
}

func newFileList(srv *drive.Service, cwd string, debug_mode bool) (*drive.FileList, error) {
	if filelist, err := api.GetFolderContent(srv, cwd); err != nil {
		return nil, api.ToHumanReadableError(err, debug_mode)
	} else {
		return ui.SortFileList(filelist), nil
	}
}

func (m browserModel) loadNewFileList() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			filelist, err := newFileList(m.srv, strings.Join(m.cwd, ""), m.debug_mode)
			return doneLoadingMsg{filelist, err}
		},
	)
}

func fileListToItems(filelist *drive.FileList) []list.Item {
	items := []list.Item{}
	for i := range filelist.Files {
		items = append(items, browserItem(fmt.Sprintf("%d", i)))
	}
	return items
}

func (m browserModel) Init() tea.Cmd {
	return m.loadNewFileList()
}

func (m browserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case doneLoadingMsg:
		m.filelist = msg.filelist
		m.list.SetDelegate(itemRenderer{m.filelist, m.format})
		m.list.SetItems(fileListToItems(m.filelist))
		m.list.Title = m.CWD()
		m.loading = false

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Root):
			if len(m.cwd) > 0 {
				m.cwd = []string{}
				m.loading = true
				return m, m.loadNewFileList()
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Backspace):
			if len(m.cwd) > 0 {
				m.cwd = m.cwd[:len(m.cwd)-1]
				m.loading = true
				return m, m.loadNewFileList()
			}

		case key.Matches(msg, m.keys.Select):
			_, ok := m.list.SelectedItem().(browserItem)
			if ok {
				file := m.filelist.Files[m.list.Index()]
				if file.MimeType == "application/vnd.google-apps.folder" {
					m.cwd = append(m.cwd, file.Name)
					m.loading = true
					return m, m.loadNewFileList()
				}
			}
			return m, nil
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m browserModel) View() string {
	if m.quitting {
		return ""
	}
	if m.loading {
		m.list.Title = fmt.Sprintf("%s loading %s...", m.spinner.View(), m.CWD())
		m.list.Styles.Title = lipgloss.NewStyle().MarginLeft(2)
	}
	return "\n" + m.list.View() + "\n" + m.help.View(m.keys)
}
