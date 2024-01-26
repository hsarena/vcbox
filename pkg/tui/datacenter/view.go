package datacenter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/muesli/reflow/wordwrap"
)

func (bd BubbleDatacenter) View(svt common.ShowViewType, height int) string {
	switch svt {
	case common.ShowList:
		return bd.listView(height)
	case common.ShowDetail:
		return bd.detailView()
	case common.ShowFull:
		return bd.fullView(height)
	default:
		return bd.fullView(height)
	}
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

func renderDCDetails(i item) string {
	dcName := fmt.Sprintf("\nName: %s", i.name)
	totalHost := fmt.Sprintf("\tHosts: %v", i.totalHosts)
	totalVMs := fmt.Sprintf("\tVMs: %v", i.totalVMs)

	return dcName + totalHost + totalVMs
}
