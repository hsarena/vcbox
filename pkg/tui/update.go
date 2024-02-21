package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		headerCmd, sideCmd, tabCmd tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.header.Height = msg.Height
		m.header.Width = msg.Width
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "down", "tab":
			m.tab, tabCmd = m.tab.Update(msg)
			m.header, headerCmd = m.header.Update(msg)
			m.side, sideCmd = m.side.Update(msg)
			return m, tea.Batch(headerCmd, sideCmd, tabCmd)
		case "left", "right":
			m.tab, tabCmd = m.tab.Update(msg)
			return m, tabCmd
		}
	}
	return m, tea.Batch(headerCmd, sideCmd, tabCmd)
}
