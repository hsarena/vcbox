package host

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/hsarena/vcbox/pkg/util"
	"github.com/muesli/reflow/wordwrap"
)

func (bd BubbleHost) View(svt common.ShowViewType, height int) string {
	switch svt {
	case common.ShowList:
		return bd.listView(height)
	case common.ShowDetail:
		//return bd.logView()
		return bd.metricsView()
	case common.ShowFull:
		return bd.fullView(height)
	default:
		return bd.fullView(height)
	}
}

func (bh BubbleHost) fullView(height int) string {
	bh.viewport.SetContent(bh.metricsView())
	//bh.viewport.SetContent(bh.logView())
	return lipgloss.JoinHorizontal(
		lipgloss.Top, bh.listView(height), bh.viewport.View())
}

func (bh BubbleHost) listView(height int) string {
	bh.list.Styles.Title = common.ListColorStyle
	bh.list.Styles.FilterPrompt.Foreground(common.ListColorStyle.GetBackground())
	bh.list.Styles.FilterCursor.Foreground(common.ListColorStyle.GetBackground())
	bh.list.SetHeight(height / 8)
	return common.ListStyle.Render(bh.list.View())
}

func (bh BubbleHost) logView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bh.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")
	logHeader := common.HeaderStyle.Render("Logs")

	if it := bh.list.SelectedItem(); it != nil {
		builder.WriteString(detailsHeader)
		builder.WriteString("\n")
		builder.WriteString(renderHostDetails(it.(item)))
		builder.WriteString(divider)
		builder.WriteString(logHeader)
		builder.WriteString("\n\n")
		if util.IsMock() {
			builder.WriteString(fmt.Sprintf("This is the log of host[%s]", it.(item).name))
		} else {
			builder.WriteString(renderHostLog(it.(item)))
		}
	}
	details := wordwrap.String(builder.String(), bh.viewport.Width)

	return common.DetailStyle.Render(details)
}

func (bh BubbleHost) metricsView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bh.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")
	metricsHeader := common.HeaderStyle.Render("Metrics")

	if it := bh.list.SelectedItem(); it != nil {
		hostMetrics, err := bh.metrics.FetchMetrics(it.(item).obj, common.HostMetrics)
		if err != nil {
			log.Printf("%s", err.Error())
		}
		graph := make([]string, len(hostMetrics))
		builder.WriteString(detailsHeader)
		builder.WriteString("\n")
		builder.WriteString(renderHostDetails(it.(item)))
		builder.WriteString(divider)
		builder.WriteString(metricsHeader)
		builder.WriteString("\n\n")
		for i, x := range hostMetrics {
			vf64 := util.ToF64(x.Value)
			graph = append(graph, asciigraph.Plot(vf64, asciigraph.SeriesColors(asciigraph.DarkGoldenrod),
				asciigraph.AxisColor(asciigraph.IndianRed),
				asciigraph.Height(bh.viewport.Height/3),
				asciigraph.Width(bh.viewport.Width/3),
				asciigraph.Caption(util.MetricIdToString(i)),
				asciigraph.Offset(2)))
		}
		builder.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, graph...))
	}
	details := wordwrap.String(builder.String(), bh.viewport.Width)

	return common.DetailStyle.Render(details)
}

func renderHostDetails(i item) string {

	hostName := fmt.Sprintf("Name: %s", i.name)
	uptime := fmt.Sprintf("\tUptime: %v days",i.uptime)
	powerState := fmt.Sprintf("\tStatus: %v",i.powerState)
	cpuModel := fmt.Sprintf("\nCPU Model: %v",i.cpuModel)
	numCpuCores := fmt.Sprintf("\tCPU Cores: %v",i.numCpuCores)
	memorySize := fmt.Sprintf("\nMemory: %vGB",i.memorySize)
	numNics := fmt.Sprintf("\tNics: %v",i.numNics)
	numHBAs := fmt.Sprintf("\tHBAs: %v",i.numHBAs)
	thostMaxVirtualDiskCapacity := fmt.Sprintf("\nhostMaxVirtualDiskCapacity: %vTB",i.hostMaxVirtualDiskCapacity)
	return hostName + uptime + powerState + cpuModel + numCpuCores + memorySize + numNics + numHBAs + thostMaxVirtualDiskCapacity
}

func renderHostLog(i item) string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	buf := new(bytes.Buffer)
	err := i.logs.Seek(ctx, 20)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	i.logs.Copy(ctx, buf)
	return buf.String()
}
