package vm

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

func (bv BubbleVM) View(svt common.ShowViewType, height int) string {
	switch svt {
	case common.ShowList:
		return bv.listView(height)
	case common.ShowDetail:
		//return bv.detailView()
		return bv.metricsView()
	case common.ShowFull:
		return bv.fullView(height)
	default:
		return bv.fullView(height)
	}
}

func (bv BubbleVM) fullView(height int) string {
	//bv.viewport.SetContent(bv.detailView())
	bv.viewport.SetContent(bv.metricsView())
	return lipgloss.JoinHorizontal(
		lipgloss.Top, bv.listView(height), bv.viewport.View())
}

func (bv BubbleVM) listView(height int) string {
	bv.list.Styles.Title = common.ListColorStyle
	bv.list.Styles.FilterPrompt.Foreground(common.ListColorStyle.GetBackground())
	bv.list.Styles.FilterCursor.Foreground(common.ListColorStyle.GetBackground())
	bv.list.SetHeight(height / 2)
	return common.ListStyle.Render(bv.list.View())
}

func (bv BubbleVM) detailView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bv.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")
	if it := bv.list.SelectedItem(); it != nil {

		builder.WriteString(detailsHeader)
		builder.WriteString("\n")
		builder.WriteString(renderVMDetails(it.(item)))
		builder.WriteString(divider)
	}
	details := wordwrap.String(builder.String(), bv.viewport.Width)

	return common.DetailStyle.Render(details)
}

func (bv BubbleVM) metricsView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bv.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")
	metricsHeader := common.HeaderStyle.Render("Metrics")

	if it := bv.list.SelectedItem(); it != nil {
		vmMetrics, err := bv.metrics.FetchMetrics(it.(item).vm, common.VMMetrics)
		if err != nil {
			log.Printf("%s", err.Error())
		}
		graph := make([]string, len(vmMetrics))
		builder.WriteString(detailsHeader)
		builder.WriteString("\n")
		builder.WriteString(renderVMDetails(it.(item)))
		builder.WriteString(divider)
		builder.WriteString(metricsHeader)
		builder.WriteString("\n\n")
		for i, x := range vmMetrics {
			vf64 := util.ToF64(x.Value)
			graph = append(graph, asciigraph.Plot(vf64, asciigraph.SeriesColors(asciigraph.DarkGoldenrod),
				asciigraph.AxisColor(asciigraph.IndianRed),
				asciigraph.Height(bv.viewport.Height/5),
				asciigraph.Width(bv.viewport.Width/7),
				asciigraph.Caption(util.MetricIdToString(i)),
				asciigraph.Offset(2)))
		}
		builder.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, graph...))
	}
	details := wordwrap.String(builder.String(), bv.viewport.Width)

	return common.DetailStyle.Render(details)
}

func renderVMDetails(i item) string {
	vmName := fmt.Sprintf("\nName: %s", i.name)
	vmOS := fmt.Sprintf("\nOS: %s", i.os)
	vmCPU := fmt.Sprintf("\nCPU: %v", i.cpu)
	vmMemory := fmt.Sprintf("\nMemory: %v", i.memory)
	vmIP := fmt.Sprintf("\nIP: %v", i.ip)
	vmStatus := fmt.Sprintf("\nStatus: %s", i.status)

	return vmName + vmOS + vmCPU + vmMemory + vmIP + vmStatus
}
