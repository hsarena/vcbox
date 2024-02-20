package side

import "github.com/charmbracelet/lipgloss"

var (
	ListStyle = lipgloss.NewStyle().
			PaddingRight(1).
			MarginRight(1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#f56642"))
	DocStyle       = lipgloss.NewStyle().Margin(1, 2)
	ListColorStyle = lipgloss.NewStyle().
			Background(lipgloss.NoColor{}).
			Foreground(lipgloss.Color("#00ffa2"))
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)
	ListSelectedListItemStyle = lipgloss.NewStyle().
					PaddingLeft(2).
					Foreground(lipgloss.Color("#00ffa2"))
	SideViewStyle = lipgloss.NewStyle().
			Padding(1)
)
