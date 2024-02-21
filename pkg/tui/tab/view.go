package tab

import (
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/hsarena/vcbox/pkg/util"
	"github.com/muesli/reflow/wordwrap"
	"github.com/vmware/govmomi/vim25/types"
)

func (m Model) View() string {

	page := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabBarStyle.Copy()
		} else {
			style = inactiveTabBarStyle.Copy()
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
		renderedTabs = append(renderedTabs, style.Render(getKindName(t.kind)))
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	page.WriteString(tabBar)
	page.WriteString("\n")
	page.WriteString(pageStyle.Render(m.tabs[m.activeTab].content))
	return wordwrap.String(page.String(), m.view.Width)

}

func getKindName(k tabKind) string {
	switch k {
	case Metrics:
		return "Metrics"
	case Logs:
		return "Logs"
	case Remote:
		return "Remote"
	default:
		return ""
	}
}

func (m Model) metricsView(obj types.ManagedObjectReference) string {
	builder := &strings.Builder{}
	metrics, err := m.metrics.FetchMetrics(obj, vmMetrics)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	graph := []string{}
	for i, x := range metrics {
		graph = append(graph, asciigraph.Plot(x, asciigraph.SeriesColors(asciigraph.DarkGoldenrod),
			asciigraph.AxisColor(asciigraph.IndianRed),
			asciigraph.Height(m.view.Height/6),
			asciigraph.Width(m.view.Width/6),
			asciigraph.Caption(util.MetricIdToString(i)),
			asciigraph.Offset(5)))

		builder.WriteString(lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.JoinHorizontal(lipgloss.Top, graph[:len(graph)/2]...), "\n\n\n",
			lipgloss.JoinHorizontal(lipgloss.Top, graph[len(graph)/2:]...)))

	}
	details := wordwrap.String(builder.String(), m.view.Width)

	return metricsStyle.Render(details)
}
