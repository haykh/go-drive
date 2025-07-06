package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	BrowserItem         = lipgloss.NewStyle().PaddingLeft(4)
	BrowserSelectedItem = WithForeground("170").PaddingLeft(2)
	BrowserPagination   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	BrowserHelp         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	BrowserDebugTitle   = WithForeground("0").PaddingLeft(2).PaddingRight(2).Background(lipgloss.Color("1"))
	BrowserMainField    = lipgloss.NewStyle().Width(90)
)
