package tab

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vmware/govmomi/vim25/types"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.view = viewport.New(msg.Width, msg.Width)
		m.view.SetContent(m.View(types.ManagedObjectReference{}))
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "right", "l", "n":
			m.activeTab = tabKind(min(int(m.activeTab)+1, len(m.tabs)-1))
			return m, nil
		case "left", "h", "p":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)

	return m, cmd
}
