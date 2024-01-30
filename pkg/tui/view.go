package tui

import (
	"strings"

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
		m.tabs = m.bd.Tab()
		//detail = m.bd.View(common.ShowMetric, m.height)
	case showHostView:
		m.tabs = m.bh.Tab()
		//detail = m.bh.View(common.ShowMetric, m.height)
	case showVMView:
		m.tabs = m.bv.Tab()
		//detail = m.bv.View(common.ShowMetric, m.height)
	default:
		m.tabs = m.bd.Tab()
		//detail = m.bd.View(common.ShowFull, m.height)
	}

	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = common.ActiveTabStyle.Copy()
		} else {
			style = common.InactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.Name))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(common.WindowStyle.Width(m.width - lipgloss.Width(side) - 4).Render(m.tabs[m.activeTab].Content))
	detail = common.DocStyle.Render(doc.String())
	// tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	// tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}
	// m := model{Tabs: tabs, TabContent: tabContent}

	return lipgloss.JoinHorizontal(lipgloss.Top, side, detail)
}
