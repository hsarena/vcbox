package tab

import "github.com/charmbracelet/lipgloss"

var (
	inactiveTabBarBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBarBorder   = tabBorderWithBottom("┘", " ", "└")
	inactiveTabBarStyle  = lipgloss.NewStyle().Border(inactiveTabBarBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabBarStyle    = inactiveTabBarStyle.Copy().Border(activeTabBarBorder, true)
	pageStyle            = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(1, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	highlightColor       = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	metricsStyle         = lipgloss.NewStyle().
				Padding(1).
				Align(lipgloss.Left)

	hostMetrics = []string{"cpu.usagemhz.average",
		"mem.consumed.average"}
	vmMetrics = []string{"cpu.usage.average",
		"mem.usage.average",
		"net.usage.average",
		"virtualDisk.write.average",
		"virtualDisk.read.average",
	}
	dcMetrics = []string{"vmop.numPoweron.latest",
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
