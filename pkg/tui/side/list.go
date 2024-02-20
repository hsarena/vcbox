package side

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

type kind int

const (
	Datacenter kind = iota
	Host
	VM
)

var (
	kinds = []kind{Datacenter, Host, VM}
)

type item struct {
	name        string
	totalVMs    int
	totalHosts  int
	cpu         int32
	memory      int32
	os          string
	ip          string
	status      string
	cpuModel    string
	numCpuCores int16
	memorySize  int64
	uptime      int32
	numNics     int32
	numHBAs     int32
	powerState  string
	logs        *object.DiagnosticLog
	obj         types.ManagedObjectReference
}

func (i item) Object() types.ManagedObjectReference { return i.obj }
func (i item) Title() string                        { return i.name }
func (i item) FilterValue() string                  { return i.name }

type itemDelegate struct{}

func (id itemDelegate) Height() int                               { return 1 }
func (id itemDelegate) Spacing() int                              { return 0 }
func (id itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (id itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	line := i.name

	if index == m.Index() {
		line = listSelectedListItemStyle.Render("_|" + line)
	} else {
		line = listItemStyle.Render(line)
	}

	fmt.Fprint(w, line)
}

func dcToItem(dcs []vmware.Inventory) []list.Item {
	items := make([]list.Item, len(dcs))
	for i, d := range dcs {
		items[i] = item{
			name:       d.Datacenter.Name(),
			totalVMs:   len(d.VMs),
			totalHosts: len(d.Hosts),
			obj:        d.Datacenter.Reference(),
		}
	}

	return items
}

func hostToItem(hosts []vmware.HostInventory) []list.Item {
	items := make([]list.Item, len(hosts))
	for i, h := range hosts {
		items[i] = item{
			name:        h.HostSystem.Name(),
			obj:         h.HostSystem.Reference(),
			logs:        h.Log,
			cpuModel:    h.CpuModel,
			numCpuCores: h.NumCpuCores,
			memorySize:  h.MemorySize / 1024 / 1024 / 1024,
			uptime:      h.Uptime / 60 / 60 / 24,
			numNics:     h.NumNics,
			numHBAs:     h.NumHBAs,
			powerState:  h.PowerState,
		}
	}

	return items
}

func vmToItem(vms []vmware.VMInventory) []list.Item {
	items := make([]list.Item, len(vms))
	for i, v := range vms {
		items[i] = item{
			name:   v.Name,
			cpu:    v.CPU,
			memory: v.Memory / 1024,
			os:     v.OS,
			ip:     v.IP,
			status: v.Status,
			obj:    v.VM.Reference(),
		}
	}
	return items
}

func getKindName(k kind) string {
	switch k {
	case Datacenter:
		return "Datacenters"
	case Host:
		return "ESXi Hosts"
	case VM:
		return "Virtaul Machines"
	default:
		return ""
	}
}

func renderItems(i item, k kind) string {
	switch k {
	case Datacenter:
		dcName := fmt.Sprintf("Name: %s", i.name)
		totalHost := fmt.Sprintf("\nHosts: %v", i.totalHosts)
		totalVMs := fmt.Sprintf("\nVMs: %v", i.totalVMs)
		return dcName + totalHost + totalVMs
	case Host:
		hostName := fmt.Sprintf("Name: %s", i.name)
		uptime := fmt.Sprintf("\tUptime: %v days", i.uptime)
		powerState := fmt.Sprintf("\tStatus: %v", i.powerState)
		cpuModel := fmt.Sprintf("\nCPU Model: %v", i.cpuModel)
		memorySize := fmt.Sprintf("\tMemory: %vGB", i.memorySize)
		numCpuCores := fmt.Sprintf("\nCPU Cores: %v", i.numCpuCores)
		numNics := fmt.Sprintf("\tNics: %v", i.numNics)
		numHBAs := fmt.Sprintf("\tHBAs: %v", i.numHBAs)
		return hostName + uptime + powerState + cpuModel + memorySize + numCpuCores + numNics + numHBAs
	case VM:
		vmName := fmt.Sprintf("Name: %s", i.name)
		vmOS := fmt.Sprintf("\tOS: %s", i.os)
		vmCPU := fmt.Sprintf("\nCPU: %v", i.cpu)
		vmMemory := fmt.Sprintf("\tMemory: %vGB", i.memory)
		vmIP := fmt.Sprintf("\nIP: %v", i.ip)
		vmStatus := fmt.Sprintf("\tStatus: %s", i.status)
		return vmName + vmOS + vmCPU + vmMemory + vmIP + vmStatus
	default:
		return i.name
	}
}
