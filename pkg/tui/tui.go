package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/datacenter"
	"github.com/hsarena/vcbox/pkg/tui/host"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi"
)

type state int

const (
	showDatacenterView state = iota
	showHostView
	showVMView
)

type model struct {
	state         state
	bd            datacenter.BubbleDatacenter
	bh            host.BubbleHost
	width, height int
	inventory     []vmware.Inventory
}

func (m model) Init() tea.Cmd {
	return nil
}

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
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {

	// view := m.bd.ListView() + "\n" + m.bh.ListView()
	// return view
	switch m.state {
	case showDatacenterView:
		return m.bd.View()
	case showHostView:
		return m.bh.View()
	// case showVMView:
	// 	return m.bd.ListView()
	default:
		return m.bd.View()
	}
}

func InitialModel(client *govmomi.Client) model {
	d := vmware.NewDiscoveryService(client)
	inventory, err := d.FetchInventory()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	return model{
		state:     showDatacenterView,
		bd:        datacenter.InitialModel(inventory),
		bh:        host.InitialModel(inventory[0].Hosts),
		inventory: inventory,
	}
}

func updateByState(m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	windowSizeMsg := tea.WindowSizeMsg{
		Width:  m.width,
		Height: m.height,
	}

	if m.state == showDatacenterView {
		m.state = showHostView
		m.bh, cmd = m.bh.Update(windowSizeMsg)
	} else {
		m.state = showDatacenterView
		m.bd, cmd = m.bd.Update(windowSizeMsg)
	}

	return m, cmd
}
