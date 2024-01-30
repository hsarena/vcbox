package vm

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi/vim25/types"
)

type item struct {
	name   string
	cpu    int32
	memory int32
	os     string
	ip     string
	status string
	vm     types.ManagedObjectReference
}

func (i item) Name() string        { return i.name }
func (i item) Description() string { return i.ip }
func (i item) FilterValue() string { return i.name }

type BubbleVM struct {
	list     list.Model
	viewport viewport.Model
	metrics  *vmware.MetricsService
}

func NewBubbleVM(l list.Model) BubbleVM {
	return BubbleVM{list: l}
}

func InitialModel(inventory []vmware.VMInventory, metrics *vmware.MetricsService) BubbleVM {
	items := vmToItem(inventory)
	l := list.New(items, vmItemDelegate{}, 0, 0)
	l.Title = "VirtualMachines"
	l.SetShowHelp(true)
	l.SetShowStatusBar(true)
	return BubbleVM{list: l, metrics: metrics}
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
			vm:     v.VM,
		}
	}
	return items
}

func MockInitialModel(inventory []vmware.MockVMInventory) BubbleVM {
	items := mockVMToItem(inventory)
	l := list.New(items, vmItemDelegate{}, 0, 0)
	l.Title = "Virtual Machines"
	l.SetShowHelp(true)
	l.SetShowStatusBar(true)
	return BubbleVM{list: l}
}

func mockVMToItem(vms []vmware.MockVMInventory) []list.Item {
	items := make([]list.Item, len(vms))
	for i, v := range vms {
		items[i] = item{
			name:   v.Name,
			cpu:    v.CPU,
			memory: v.Memory,
			os:     v.OS,
			ip:     v.IP,
			status: v.Status,
		}
	}
	return items
}
