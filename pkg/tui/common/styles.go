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
	MetricsStyle = lipgloss.NewStyle().
			Padding(1).
			Align(lipgloss.Left)

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
