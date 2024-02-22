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

func (m Model) View(selected types.ManagedObjectReference) string {

	page := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, tabKind(i) == m.activeTab
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

	switch m.activeTab {
	case Metrics:
		m.tabs[m.activeTab].content = m.metricsView(selected)
	}

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
			asciigraph.Height(5),
			asciigraph.Width(20),
			asciigraph.Caption(util.MetricIdToString(i)),
		))
	}
	totalGraphs := len(graph)
	numRows := totalGraphs / 2
	graphsPerRow := 2

	// Calculate the total number of graphs

	// Iterate over each row
	for i := 0; i < numRows; i++ {
		// Calculate the start and end indices for the graphs in the current row
		start := i * graphsPerRow
		end := start + graphsPerRow
		if i == numRows-1 {
			// If it's the last row, adjust the end index to include the remaining graphs
			end = totalGraphs
		}
		// Join the graphs for the current row horizontally
		row := lipgloss.JoinHorizontal(lipgloss.Top, graph[start:end]...)
		// Append the current row to the string builder
		builder.WriteString(row)
		// Add a newline character to separate rows
		builder.WriteString("\n")
	}
	details := wordwrap.String(builder.String(), m.view.Width)

	return metricsStyle.Render(details)
}
