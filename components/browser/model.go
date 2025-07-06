package browser

import (
	"bytes"
	"fmt"
	"go-drive/filesystem"
	"go-drive/ui"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type browserModel struct {
	// model configurations
	debug_mode bool
	format     []string

	// model state
	cwd         []string
	filemanager filesystem.FileManager
	filelist    []filesystem.FileItem
	filesinsync []int

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
		m.list.SetDelegate(itemRenderer{m.filelist, m.filesinsync, m.cwd})
		m.list.SetItems(FileListToItems(m.filelist))
		m.list.Title = m.CWD()
		if m.list.Title == "" {
			m.list.Title = "/"
		}
		m.loading = false

	case syncFileMsg:
		if !slices.Contains(m.filesinsync, msg.index) {
			m.filesinsync = append(m.filesinsync, msg.index)
		}
		m.list.SetDelegate(itemRenderer{m.filelist, m.filesinsync, m.cwd})

	case doneSyncMsg:
		m.filelist = msg.filelist
		for i, idx := range m.filesinsync {
			if idx == msg.index {
				m.filesinsync = append(m.filesinsync[:i], m.filesinsync[i+1:]...)
				break
			}
		}
		m.list.SetDelegate(itemRenderer{m.filelist, m.filesinsync, m.cwd})
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
			if len(m.debugLines) > 30 {
				m.debugLines = m.debugLines[len(m.debugLines)-30:]
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
				return m.changeDirectory([]string{})
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Backspace):
			if len(m.cwd) > 0 {
				return m.changeDirectory(m.cwd[:len(m.cwd)-1])
			}

		case key.Matches(msg, m.keys.Select):
			if _, ok := m.list.SelectedItem().(browserItem); ok {
				if file := m.filelist[m.list.Index()]; file.IsDirectory() {
					return m.changeDirectory(append(m.cwd, file.GetName()))
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Sync):
			if _, ok := m.list.SelectedItem().(browserItem); ok {
				file := m.filelist[m.list.Index()]
				if !file.IsDirectory() {
					m.loading = true
					m.status = fmt.Sprintf("syncing %s...", file.GetName())
					return m, m.syncFile(m.list.Index(), file)
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
		debug_log = fmt.Sprintf("%s\n%s", ui.BrowserDebugTitle.Render("debug log"), strings.Join(m.debugLines, "\n"))
	}
	normal_view := fmt.Sprintf("\n%s\n%s\n", m.list.View(), ui.BrowserHelp.Render(m.help.View(m.keys)))
	return lipgloss.JoinHorizontal(lipgloss.Top, ui.BrowserMainField.Render(normal_view), debug_log)
}
