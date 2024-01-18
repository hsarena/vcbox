package tui

import tea "github.com/charmbracelet/bubbletea"

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
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "tab" {
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
		m.bh, cmd = m.bh.Update(windowSizeMsg)
	case showHostView:
		m.state = showVMView
		m.bd, cmd = m.bd.Update(windowSizeMsg)
	case showVMView:
		m.state = showDatacenterView
		m.bv, cmd = m.bv.Update(windowSizeMsg)
	default:
		m.state = showDatacenterView
		m.bd, cmd = m.bd.Update(windowSizeMsg)
	}

	return m, cmd
}
