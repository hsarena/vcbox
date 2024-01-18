package host

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/muesli/reflow/wordwrap"
)

func (bh BubbleHost) View() string {
	bh.viewport.SetContent(bh.detailView())

	return lipgloss.JoinHorizontal(
		lipgloss.Top, bh.listView(), bh.viewport.View())
}

func (bh BubbleHost) listView() string {
	bh.list.Styles.Title = common.ListColorStyle
	bh.list.Styles.FilterPrompt.Foreground(common.ListColorStyle.GetBackground())
	bh.list.Styles.FilterCursor.Foreground(common.ListColorStyle.GetBackground())

	return common.ListStyle.Render(bh.list.View())
}

func (bh BubbleHost) detailView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bh.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Details")
	if it := bh.list.SelectedItem(); it != nil {

		builder.WriteString(detailsHeader)
		builder.WriteString(renderHostDetails(it.(item)))
		builder.WriteString(divider)
		builder.WriteString(renderHostLog(it.(item)))
	}
	details := wordwrap.String(builder.String(), bh.viewport.Width)

	return common.DetailStyle.Render(details)
}

func (bh BubbleHost) logView() string {
	builder := &strings.Builder{}
	divider := common.DividerStyle.Render(strings.Repeat("-", bh.viewport.Width)) + "\n"
	detailsHeader := common.HeaderStyle.Render("Logs")

	if it := bh.list.SelectedItem(); it != nil {
		builder.WriteString(detailsHeader)
		builder.WriteString(renderHostDetails(it.(item)))
		builder.WriteString(divider)
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
	err := i.logs.Seek(ctx, 100)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	i.logs.Copy(ctx, buf)
	return buf.String()
}
