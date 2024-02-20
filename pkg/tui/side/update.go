package side

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			return m.updateByActiveSide(msg)
		}
	}
	var cmd tea.Cmd
	m.sides[m.activeSide].list, cmd = m.sides[m.activeSide].list.Update(msg)
	return m, cmd
}

func (m Model) updateByActiveSide(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	selected := m.sides[Datacenter].list.Index()
	switch m.activeSide {
	case Datacenter:
		m.activeSide = Host
		m.sides[Datacenter].list, cmd = m.sides[Datacenter].list.Update(msg)
		return m, cmd
	case Host:
		m.activeSide = VM
		return m, m.sides[Host].list.SetItems(hostToItem(m.inventory[selected].Hosts))
	case VM:
		m.activeSide = Datacenter
		return m, m.sides[VM].list.SetItems(vmToItem(m.inventory[selected].VMs))
	default:
		m.activeSide = Datacenter
		m.sides[Datacenter].list, cmd = m.sides[Datacenter].list.Update(msg)
		return m, cmd
	}
}
