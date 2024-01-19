package mock

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
	in := make([]vmware.MockInventory, len(dcD))
	for i, dc := range dcD {
		in[i].Datacenter = dc
		crD := [3]string{"srv-01", "srv-02", "srv-03"}
		hostI := make([]vmware.MockHostInventory, len(crD))
		for i, c := range crD {
			hostI[i].ComputeResource = c
			hostI[i].Log = fmt.Sprintf("This is the mock log of host [%s]", c)
		}
		vms := [10]string{"VM-01",
			"VM-02",
			"VM-03",
			"VM-04",
			"VM-05",
			"VM-06",
			"VM-07",
			"VM-08",
			"VM-09",
			"VM-10"}

		vmI := make([]vmware.MockVMInventory, len(vms))
		for i := range vms {
			vm := &vmware.MockVMInventory{
				Name:   vms[i],
				OS:     "openSUSE",
				CPU:    4,
				Memory: 4096,
				IP:     "192.168.10.100",
				Status: "On",
			}
			vmI[i] = *vm
		}
		in[i].Hosts = hostI
		in[i].VMs = vmI
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
		m.bh, cmd = m.bh.Update(windowSizeMsg)
	case showHostView:
		m.state = showVMView
		//m.bd, cmd = m.bd.Update(windowSizeMsg)
		m.bv, cmd = m.bv.Update(windowSizeMsg)
	case showVMView:
		m.state = showDatacenterView
		m.bv, cmd = m.bv.Update(windowSizeMsg)
	default:
		m.state = showDatacenterView
		m.bd, cmd = m.bd.Update(windowSizeMsg)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case showDatacenterView:
		return m.bd.View(true)
	case showHostView:
		return m.bh.View(true)
	case showVMView:
		return m.bv.View(true)
	default:
		return m.bd.View(true)
	}
}
