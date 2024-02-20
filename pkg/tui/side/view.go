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
		l.Styles.Title = ListColorStyle
		l.Title = getKindName(k)
		
		l.Styles.FilterPrompt.Foreground(ListColorStyle.GetBackground())
		l.Styles.FilterCursor.Foreground(ListColorStyle.GetBackground())
		renderedList = append(renderedList, ListStyle.Render(l.View()))
	}
	m.view.Style = SideViewStyle
	m.view.SetContent(lipgloss.JoinVertical(lipgloss.Top, renderedList...))
	return m.view.View()
}
