package browser

import (
	"bytes"
	"fmt"
	"go-drive/utils"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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

type errorMsg struct {
	err error
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
	loading     bool
	status      string
	quitting    bool
	debugLines  []string
	debugBuffer *bytes.Buffer

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

func (m browserModel) syncFile(file utils.FileItem) tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			mgr := m.filemanager
			if err := mgr.Synchronize(file, m.debug_mode); err != nil {
				return errorMsg{
					err: fmt.Errorf("failed to synchronize file %s: %w", file.GetName(), err),
				}
			}
			filelist, err := mgr.GetFileList(m.CWD(), m.debug_mode)
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

	case errorMsg:
		m.loading = false
		log.Errorf("error loading file list: %s", msg.err.Error())

	case spinner.TickMsg:
		if m.debug_mode {
			data := m.debugBuffer.String()
			m.debugBuffer.Reset()
			debug_lines := []string{}
			if data != "" {
				debug_lines = strings.Split(strings.TrimSuffix(data, "\n"), "\n")
			}
			m.debugLines = append(m.debugLines, debug_lines...)
			if len(m.debugLines) > 40 {
				m.debugLines = m.debugLines[len(m.debugLines)-40:]
			}
		}
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
				m.status = "loading /..."
				return m, m.loadRemoteFileList()
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Backspace):
			if len(m.cwd) > 0 {
				m.cwd = m.cwd[:len(m.cwd)-1]
				m.loading = true
				m.status = fmt.Sprintf("loading %s...", m.CWD())
				return m, m.loadRemoteFileList()
			}

		case key.Matches(msg, m.keys.Select):
			if _, ok := m.list.SelectedItem().(browserItem); ok {
				file := m.filelist[m.list.Index()]
				if file.IsDirectory() {
					m.cwd = append(m.cwd, file.GetName())
					m.loading = true
					m.status = fmt.Sprintf("loading %s...", m.CWD())
					return m, m.loadRemoteFileList()
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Sync):
			if _, ok := m.list.SelectedItem().(browserItem); ok {
				file := m.filelist[m.list.Index()]
				if !file.IsDirectory() {
					m.loading = true
					m.status = fmt.Sprintf("syncing %s...", file.GetName())
					return m, m.syncFile(file)
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
		m.list.Title = fmt.Sprintf("%s %s", m.spinner.View(), m.status)
		m.list.Styles.Title = lipgloss.NewStyle().MarginLeft(2)
	}
	debug_log := ""
	if m.debug_mode {
		debug_log = fmt.Sprintf("%s\n%s", debugTitle.Render("debug log"), strings.Join(m.debugLines, "\n"))
	}
	normal_view := fmt.Sprintf("\n%s\n%s\n", m.list.View(), helpStyle.Render(m.help.View(m.keys)))
	return lipgloss.JoinHorizontal(lipgloss.Top, mainField.Render(normal_view), debug_log)
}
