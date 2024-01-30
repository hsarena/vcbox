package common

import "github.com/charmbracelet/lipgloss"

type ShowViewType int

const (
	ShowList ShowViewType = iota
	ShowDetail
	ShowMetric
	ShowLog
	ShowFull
)

var (
	ListStyle = lipgloss.NewStyle().
			Width(30).
			PaddingRight(1).
			MarginRight(1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#f56642"))
	ListColorStyle = lipgloss.NewStyle().
			Background(lipgloss.NoColor{}).
			Foreground(lipgloss.Color("#00ffa2"))
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)
	ListSelectedListItemStyle = lipgloss.NewStyle().
					PaddingLeft(2).
					Foreground(lipgloss.Color("#00ffa2"))
	DetailStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(1)
	MetricsStyle = lipgloss.NewStyle().
			Padding(1).
			Align(lipgloss.Left)
	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f56642")).
			PaddingTop(1)
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffa2")).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#f56642")).
			Bold(true)

	InactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = tabBorderWithBottom("┘", " ", "└")
	DocStyle          = lipgloss.NewStyle().Padding(1)
	HighlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	InactiveTabStyle  = lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderForeground(HighlightColor).Padding(0, 1)
	ActiveTabStyle    = InactiveTabStyle.Copy().Border(ActiveTabBorder, true)
	WindowStyle       = lipgloss.NewStyle().BorderForeground(HighlightColor).Padding(1, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()

	HostMetrics = []string{"cpu.usagemhz.average",
		"mem.consumed.average"}
	VMMetrics = []string{"cpu.usage.average",
		"mem.usage.average",
		"net.usage.average",
		"virtualDisk.write.average",
		"virtualDisk.read.average",
	}
	DCMetrics = []string{"vmop.numPoweron.latest",
		"vmop.numPoweroff.latest",
		"vmop.numCreate.latest",
		"vmop.numReconfigure.latest",
		"vmop.numVMotion.latest",
	}
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
