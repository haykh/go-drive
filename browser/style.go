package browser

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	debugTitle        = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).Background(lipgloss.Color("1")).Foreground(lipgloss.Color("0"))
	mainField         = lipgloss.NewStyle().Width(90)
)
