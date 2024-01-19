package vm

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/muesli/reflow/wordwrap"
)

func (bv BubbleVM) View(showDetails bool) string {
	if showDetails {
		bv.viewport.SetContent(bv.detailView())
		return lipgloss.JoinHorizontal(
			lipgloss.Top, bv.listView(), bv.viewport.View())
	} else {
		return bv.listView()
	}

}

func (bv BubbleVM) listView() string {
	bv.list.Styles.Title = common.ListColorStyle
	bv.list.Styles.FilterPrompt.Foreground(common.ListColorStyle.GetBackground())
	bv.list.Styles.FilterCursor.Foreground(common.ListColorStyle.GetBackground())

	return common.ListStyle.Render(bv.list.View())
}

func (bv BubbleVM) detailView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bv.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")
	if it := bv.list.SelectedItem(); it != nil {

		builder.WriteString(detailsHeader)
		builder.WriteString(renderVMDetails(it.(item)))
		builder.WriteString(divider)
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
