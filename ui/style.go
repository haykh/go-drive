package ui

import "github.com/charmbracelet/lipgloss"

func WithColor(c string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(c))
}
