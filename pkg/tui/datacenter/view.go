package datacenter

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/hsarena/vcbox/pkg/util"
	"github.com/muesli/reflow/wordwrap"
)

func (bd BubbleDatacenter) View(svt common.ShowViewType, height int) string {
	switch svt {
	case common.ShowList:
		return bd.listView(height)
	case common.ShowDetail:
		return bd.detailView()
	case common.ShowLog:
		return ""
	case common.ShowMetric:
		return bd.metricsView()
	case common.ShowFull:
		return bd.fullView(height)
	default:
		return bd.fullView(height)
	}
}

func (bd BubbleDatacenter) Tab() []common.Tab {
	tab := make([]common.Tab, 2)
	tab[0].Name = "details"
	tab[0].Content = bd.detailView()
	tab[1].Name = "metrics"
	tab[1].Content = bd.metricsView()
	return tab
}

func (bd BubbleDatacenter) fullView(height int) string {
	bd.viewport.SetContent(bd.detailView())
	return lipgloss.JoinHorizontal(
		lipgloss.Top, bd.listView(height), bd.viewport.View())
}

func (bd BubbleDatacenter) listView(height int) string {
	bd.list.Styles.Title = common.ListColorStyle
	bd.list.Styles.FilterPrompt.Foreground(common.ListColorStyle.GetBackground())
	bd.list.Styles.FilterCursor.Foreground(common.ListColorStyle.GetBackground())
	bd.list.SetHeight(height / 8)
	return common.ListStyle.Render(bd.list.View())
}

func (bd BubbleDatacenter) detailView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bd.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")

	if it := bd.list.SelectedItem(); it != nil {
		builder.WriteString(detailsHeader)
		builder.WriteString("\n")
		builder.WriteString(renderDCDetails(it.(item)))
		builder.WriteString(divider)
	}
	details := wordwrap.String(builder.String(), bd.viewport.Width)

	return common.DetailStyle.Render(details)
}

func (bd BubbleDatacenter) metricsView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bd.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")
	metricsHeader := common.HeaderStyle.Render("Metrics")

	if it := bd.list.SelectedItem(); it != nil {
		vmMetrics, err := bd.metrics.FetchMetrics(it.(item).obj, common.DCMetrics)
		if err != nil {
			log.Printf("%s", err.Error())
		}
		graph := []string{}
		builder.WriteString(detailsHeader)
		builder.WriteString("\n")
		builder.WriteString(renderDCDetails(it.(item)))
		builder.WriteString(divider)
		builder.WriteString(metricsHeader)
		builder.WriteString("\n\n")
		for i, x := range vmMetrics {
			graph = append(graph, asciigraph.Plot(x, asciigraph.SeriesColors(asciigraph.DarkGoldenrod),
				asciigraph.AxisColor(asciigraph.IndianRed),
				asciigraph.Height(bd.viewport.Height/6),
				asciigraph.Width(bd.viewport.Width/6),
				asciigraph.Caption(util.MetricIdToString(i)),
				asciigraph.Offset(5)))
		}

		builder.WriteString(lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.JoinHorizontal(lipgloss.Top, graph[:len(graph)/2]...), "\n\n\n",
			lipgloss.JoinHorizontal(lipgloss.Top, graph[len(graph)/2:]...)))

	}
	details := wordwrap.String(builder.String(), bd.viewport.Width)

	return common.MetricsStyle.Render(details)
}

func renderDCDetails(i item) string {
	dcName := fmt.Sprintf("Name: %s", i.name)
	totalHost := fmt.Sprintf("\tHosts: %v", i.totalHosts)
	totalVMs := fmt.Sprintf("\tVMs: %v", i.totalVMs)

	return dcName + totalHost + totalVMs
}
