package datacenter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/muesli/reflow/wordwrap"
)

func (bd BubbleDatacenter) View() string {
	bd.viewport.SetContent(bd.detailView())

	return lipgloss.JoinHorizontal(
		lipgloss.Top, bd.listView(), bd.viewport.View())
}

func (bd BubbleDatacenter) listView() string {
	bd.list.Styles.Title = common.ListColorStyle
	bd.list.Styles.FilterPrompt.Foreground(common.ListColorStyle.GetBackground())
	bd.list.Styles.FilterCursor.Foreground(common.ListColorStyle.GetBackground())

	return common.ListStyle.Render(bd.list.View())
}

func (bd BubbleDatacenter) detailView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bd.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")

	if it := bd.list.SelectedItem(); it != nil {
		builder.WriteString(detailsHeader)
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
