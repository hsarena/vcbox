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
	tabs := wordwrap.String(m.tab.View(), m.header.Width)
	page := lipgloss.JoinVertical(lipgloss.Top, details, tabs)
	return lipgloss.JoinHorizontal(lipgloss.Top, m.side.View(), page)
}
