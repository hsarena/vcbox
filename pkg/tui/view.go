package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
)

func (m model) View() string {
	var detail string
	side := lipgloss.JoinVertical(lipgloss.Top,
		m.bd.View(common.ShowList, m.height),
		m.bh.View(common.ShowList, m.height),
		m.bv.View(common.ShowList, m.height))
	switch m.state {

	case showDatacenterView:
		detail = m.bd.View(common.ShowMetric, m.height)
	case showHostView:
		detail = m.bh.View(common.ShowMetric, m.height)
	case showVMView:
		detail = m.bv.View(common.ShowMetric, m.height)
	default:
		detail = m.bd.View(common.ShowFull, m.height)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, side, "    ", detail)
}
