package browser

import (
	"fmt"
	"go-drive/filesystem"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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

func (m browserModel) syncFile(index int, file filesystem.FileItem) tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			return syncStartMsg{index: index}
		},
		func() tea.Msg {
			mgr := m.filemanager
			if err := mgr.Synchronize(file, m.CWD(), m.debug_mode); err != nil {
				return syncFailedMsg{
					index: index,
					err:   fmt.Errorf("failed to synchronize file %s: %w", file.GetName(), err),
				}
			}
			filelist, err := mgr.GetFileList(m.CWD(), m.debug_mode)
			return syncDoneMsg{filelist, index, err}
		},
	)
}

func (m browserModel) trashFile(index int, file filesystem.FileItem) tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			return trashStartMsg{index: index}
		},
		func() tea.Msg {
			mgr := m.filemanager
			if err := mgr.Trash(file, m.CWD(), m.debug_mode); err != nil {
				return trashFailedMsg{
					index: index,
					err:   fmt.Errorf("failed to trash file %s: %w", file.GetName(), err),
				}
			}
			filelist, err := mgr.GetFileList(m.CWD(), m.debug_mode)
			return trashDoneMsg{filelist, index, err}
		},
	)
}

func (m browserModel) changeDirectory(path []string) (tea.Model, tea.Cmd) {
	m.cwd = path
	m.loading = true
	dir := "/"
	if m.CWD() != "" {
		dir = m.CWD()
	}
	m.status = fmt.Sprintf("loading %s...", dir)
	return m, m.loadRemoteFileList()
}
