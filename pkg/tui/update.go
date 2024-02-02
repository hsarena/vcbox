package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	showDatacenterView state = iota
	showHostView
	showVMView
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n":
			log.Printf("about to turn right: %v, %v", m.activeTab, len(m.tabs))
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		case "tab":
			return updateByState(m)
		}
	}

	switch m.state {
	case showDatacenterView:
		m.bd, cmd = m.bd.Update(msg)
		return m, cmd
	case showHostView:
		m.bh, cmd = m.bh.Update(msg)
		return m, cmd
	case showVMView:
		m.bv, cmd = m.bv.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func updateByState(m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	windowSizeMsg := tea.WindowSizeMsg{
		Width:  m.width,
		Height: m.height,
	}

	switch m.state {
	case showDatacenterView:
		m.state = showHostView
		m.bh, cmd = m.bh.UpdateList(m.inventory[m.bd.GetSelectedItem()].Hosts, windowSizeMsg)
	case showHostView:
		m.state = showVMView
		m.bv, cmd = m.bv.UpdateList(m.inventory[m.bd.GetSelectedItem()].VMs, windowSizeMsg)
	case showVMView:
		m.state = showDatacenterView
		m.bv, cmd = m.bv.UpdateList(m.inventory[m.bd.GetSelectedItem()].VMs, windowSizeMsg)
	default:
		m.state = showDatacenterView
		m.bd, cmd = m.bd.Update(windowSizeMsg)
	}

	return m, cmd
}
