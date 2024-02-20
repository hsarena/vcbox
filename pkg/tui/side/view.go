package side

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var renderedList []string
	for _, k := range kinds {
		l := m.sides[k].list
		l.SetShowFilter(false)
		l.SetShowHelp(false)
		l.Styles.Title = listColorStyle
		l.Title = getKindName(k)

		l.Styles.FilterPrompt.Foreground(listColorStyle.GetBackground())
		l.Styles.FilterCursor.Foreground(listColorStyle.GetBackground())
		renderedList = append(renderedList, listStyle.Render(l.View()))
	}
	return lipgloss.JoinVertical(lipgloss.Top, renderedList...)
}

func (m Model) DetailView() string {
	return renderItems(m.sides[m.activeSide].list.SelectedItem().(item), m.activeSide)
}
