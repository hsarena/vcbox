package mock

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/hsarena/vcbox/pkg/tui/datacenter"
	"github.com/hsarena/vcbox/pkg/tui/host"
	"github.com/hsarena/vcbox/pkg/tui/vm"
	"github.com/hsarena/vcbox/pkg/vmware"
)

const (
	showDatacenterView state = iota
	showHostView
	showVMView
)

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) FetchInventory() ([]vmware.MockInventory, error) {
	dcD := [3]string{"HWB", "AST", "AFR"}
	crD := [3][]string{
		{"srv-hwb-01", "srv-hwb-02"},
		{"srv-ast-01", "srv-ast-02", "srv-ast-03"},
		{"srv-afr-01", "srv-afr-02", "srv-afr-03"}}
	vms := [3][]string{
		{"H-VM-01", "H-VM-02", "H-VM-03", "H-VM-04", "H-VM-05"},
		{"VM-01", "VM-02", "VM-03", "VM-04", "VM-05",
			"VM-06", "VM-07", "VM-08", "VM-09", "VM-10",
			"VM-11", "VM-12", "VM-13", "VM-14", "VM-15"},
		{"A-VM-01", "A-VM-02", "A-VM-03", "A-VM-04", "A-VM-05",
			"A-VM-06", "A-VM-07", "A-VM-08", "A-VM-09", "A-VM-10"}}
	in := make([]vmware.MockInventory, len(dcD))
	for id, dc := range dcD {
		in[id].Datacenter = dc
		hostI := make([]vmware.MockHostInventory, len(crD[id]))
		for i, c := range crD[id] {
			hostI[i].ComputeResource = c
			hostI[i].Log = fmt.Sprintf("This is the mock log of host [%s]", c)
		}
		vmI := make([]vmware.MockVMInventory, len(vms[id]))
		for i, v := range vms[id] {
			vm := &vmware.MockVMInventory{
				Name:   v,
				OS:     "openSUSE",
				CPU:    4,
				Memory: 4096,
				IP:     "192.168.10.100",
				Status: "On",
			}
			vmI[i] = *vm
		}
		in[id].Hosts = hostI
		in[id].VMs = vmI
	}
	return in, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel() model {
	s := NewMockService()
	inventory, _ := s.FetchInventory()
	return model{
		state:        0,
		bd:           datacenter.MockInitialModel(inventory),
		bh:           host.MockInitialModel(inventory[0].Hosts),
		bv:           vm.MockInitialModel(inventory[0].VMs),
		selectedDC:   0,
		selectedHost: 0,
		selectedVM:   0,
		inventory:    inventory,
	}
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
		m.bh, cmd = m.bh.MockUpdateList(m.inventory[m.bd.GetSelectedItem()].Hosts, windowSizeMsg)
	case showHostView:
		m.state = showVMView
		m.bv, cmd = m.bv.MockUpdateList(m.inventory[m.bd.GetSelectedItem()].VMs, windowSizeMsg)
	case showVMView:
		m.state = showDatacenterView
		m.bv, cmd = m.bv.MockUpdateList(m.inventory[m.bd.GetSelectedItem()].VMs, windowSizeMsg)
	default:
		m.state = showDatacenterView
		m.bd, cmd = m.bd.Update(windowSizeMsg)
	}

	return m, cmd
}

func (m model) View() string {
	var detail string
	side := lipgloss.JoinVertical(lipgloss.Top,
		m.bd.View(common.ShowList, m.height),
		m.bh.View(common.ShowList, m.height),
		m.bv.View(common.ShowList, m.height))
	// side := m.bd.View(common.ShowList, m.height)
	switch m.state {

	case showDatacenterView:
		detail = m.bd.View(common.ShowDetail, m.height)
	case showHostView:
		detail = m.bh.View(common.ShowDetail, m.height)
	case showVMView:
		detail = m.bv.View(common.ShowDetail, m.height)
	default:
		detail = m.bd.View(common.ShowFull, m.height)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, side, detail)
	//detail = m.bd.View(false) + m.bh.View(false) + m.bv.View(false)
}
