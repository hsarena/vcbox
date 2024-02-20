package side

import "github.com/charmbracelet/lipgloss"

var (
	listStyle = lipgloss.NewStyle().
			PaddingRight(1).
			MarginRight(1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#f56642")).
			MaxWidth(40).
			Width(35)
	listColorStyle = lipgloss.NewStyle().
			Background(lipgloss.NoColor{}).
			Foreground(lipgloss.Color("#00ffa2"))
	listItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)
	listSelectedListItemStyle = lipgloss.NewStyle().
					PaddingLeft(2).
					Foreground(lipgloss.Color("#00ffa2"))
)
