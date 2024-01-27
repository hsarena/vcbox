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
		log.Printf("about to fetch metrics for obj %v", it.(item).obj)
		hostMetrics, err := bh.metrics.FetchMetrics(it.(item).obj, common.HostMetrics)
		log.Printf("the hostMetrics is: %v", hostMetrics["cpu.usagemhz.average"].Value)
		vf64 := util.ToF64(hostMetrics["cpu.usagemhz.average"].Value)
		log.Println("about to draw a graph")
		graph := asciigraph.Plot(vf64, asciigraph.AxisColor(asciigraph.DarkGoldenrod),
			asciigraph.AxisColor(asciigraph.IndianRed),
			asciigraph.Height(10),
			asciigraph.Width(30),
			asciigraph.Caption("CPU Usagemhz Average"))

		if err != nil {
			log.Printf("%s", err.Error())
		}
		builder.WriteString(detailsHeader)
		builder.WriteString("\n")
		builder.WriteString(renderHostDetails(it.(item)))
		builder.WriteString(divider)
		builder.WriteString(metricsHeader)
		builder.WriteString("\n\n")
		builder.WriteString(graph)
	}
	details := wordwrap.String(builder.String(), bh.viewport.Width)

	return common.DetailStyle.Render(details)
}

func renderHostDetails(i item) string {
	hostName := fmt.Sprintf("\nName: %s", i.name)
	return hostName
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
