package browser

import (
	"fmt"
	"go-drive/utils"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	Remote Mode = iota
	Local
)

func FileListToItems(r []utils.FileItem) []list.Item {
	items := []list.Item{}
	for i := range r {
		items = append(items, browserItem(fmt.Sprintf("%d", i)))
	}
	return items
}

type browserItem string

func (i browserItem) FilterValue() string { return "" }

type doneLoadingMsg struct {
	filelist []utils.FileItem
	err      error
}

type browserModel struct {
	// model configurations
	debug_mode bool
	format     []string

	// model state
	cwd         []string
	filemanager utils.FileManager
	filelist    []utils.FileItem

	// ui components
	help    help.Model
	list    list.Model
	spinner spinner.Model

	// ui state
	loading  bool
	quitting bool

	// ui configurations
	keys keyMap
}

func (m browserModel) CWD() string {
	return strings.Join(m.cwd, "/")
}

func (m browserModel) loadRemoteFileList() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			filelist, err := m.filemanager.GetFileList(m.CWD(), m.debug_mode)
			return doneLoadingMsg{filelist, err}
		},
	)
}

func (m browserModel) Init() tea.Cmd {
	return m.loadRemoteFileList()
}

func (m browserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case doneLoadingMsg:
		m.filelist = msg.filelist
		m.list.SetDelegate(itemRenderer{m.filelist, m.cwd})
		m.list.SetItems(FileListToItems(m.filelist))
		m.list.Title = m.CWD()
		if m.list.Title == "" {
			m.list.Title = "/"
		}
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
				return m, m.loadRemoteFileList()
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Backspace):
			if len(m.cwd) > 0 {
				m.cwd = m.cwd[:len(m.cwd)-1]
				m.loading = true
				return m, m.loadRemoteFileList()
			}

		case key.Matches(msg, m.keys.Select):
			_, ok := m.list.SelectedItem().(browserItem)
			if ok {
				file := m.filelist[m.list.Index()]
				if file.IsDirectory() {
					m.cwd = append(m.cwd, file.GetName())
					m.loading = true
					return m, m.loadRemoteFileList()
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
	return "\n" + m.list.View() + "\n" + helpStyle.Render(m.help.View(m.keys))
}
