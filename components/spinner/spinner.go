package spinner

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func RunWithSpinner(f func() (any, error), msg, errmsg, final string, debug_mode bool) (any, error) {
	if debug_mode {
		return f()
	} else {
		s := spinner.New()
		s.Spinner = spinner.Dot

		p := tea.NewProgram(spinnerModel{
			spinner: s,
			msg:     msg,
			errmsg:  errmsg,
			final:   final,
			loading: true,
			runner:  f,
		})
		if m, err := p.Run(); err != nil {
			return nil, err
		} else {
			sm := m.(spinnerModel)
			return sm.result, sm.err
		}
	}
}

type doneMsg struct {
	result any
	err    error
}

type spinnerModel struct {
	spinner spinner.Model
	msg     string
	errmsg  string
	final   string
	result  any
	loading bool
	done    bool
	err     error
	runner  func() (any, error)
}

func runWithSpinner(f func() (any, error)) tea.Cmd {
	return func() tea.Msg {
		result, err := f()
		return doneMsg{result, err}
	}
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		runWithSpinner(m.runner),
	)
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case doneMsg:
		m.loading = false
		m.done = true
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit
	}

	return m, nil
}

func (m spinnerModel) View() string {
	if m.loading {
		return fmt.Sprintf("%s %s\n", m.spinner.View(), m.msg)
	}
	if m.err != nil {
		return fmt.Sprintf("%s : %v\n", m.errmsg, m.err)
	}
	if m.final != "" {
		return fmt.Sprintf("%s\n", m.final)
	} else {
		return ""
	}
}
