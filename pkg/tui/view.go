package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func (m Model) View() string {
	builder := &strings.Builder{}
	divider := dividerStyle.Render(strings.Repeat("-", m.header.Width)) + "\n"
	builder.WriteString(detailStyle.Render(m.side.DetailView()))
	builder.WriteString(divider)
	details := wordwrap.String(builder.String(), m.header.Width)
	return lipgloss.JoinHorizontal(lipgloss.Top, m.side.View(), details)
}
