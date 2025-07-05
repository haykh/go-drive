package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Prompt(prompt, placeholder string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()

	p := tea.NewProgram(promptModel{
		textInput: ti,
		err:       nil,
		prompt:    prompt,
	})

	if m, err := p.Run(); err != nil {
		return "", err
	} else {
		return m.(promptModel).textInput.Value(), nil
	}
}

type (
	errMsg error
)

type promptModel struct {
	textInput textinput.Model
	prompt    string
	err       error
}

func (m promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m promptModel) View() string {
	return fmt.Sprintf(
		"%s%s\n%s\n",
		WithColor("2").Render(m.prompt+": "),
		m.textInput.View()[2:],
		WithColor("8").Render("(esc to quit)"),
	)
}
