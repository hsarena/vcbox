package tui

import "github.com/charmbracelet/lipgloss"

var (
	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder())
	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f56642")).
			PaddingTop(1)
)
