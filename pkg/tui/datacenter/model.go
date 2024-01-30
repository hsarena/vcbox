package datacenter

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi/vim25/types"
)

type item struct {
	name       string
	totalVMs   int
	totalHosts int
	obj        types.ManagedObjectReference
}

func (i item) Name() string        { return i.name }
func (i item) TotalHosts() int     { return i.totalHosts }
func (i item) TotalVMs() int       { return i.totalVMs }
func (i item) FilterValue() string { return i.name }

type BubbleDatacenter struct {
	list     list.Model
	viewport viewport.Model
	metrics  *vmware.MetricsService
}

func NewBubbleDatacenter(l list.Model) BubbleDatacenter {
	return BubbleDatacenter{list: l}
}

func InitialModel(inventory []vmware.Inventory, metrics *vmware.MetricsService) BubbleDatacenter {
	items := dcToItem(inventory)
	l := list.New(items, dcItemDelegate{}, 0, 0)
	l.Title = "Datacenters"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return BubbleDatacenter{list: l, metrics: metrics}
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

func MockInitialModel(inventory []vmware.MockInventory) BubbleDatacenter {
	items := mockDCToItem(inventory)
	l := list.New(items, dcItemDelegate{}, 0, 0)
	l.Title = "Datacenters"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return BubbleDatacenter{list: l}
}

func mockDCToItem(dcs []vmware.MockInventory) []list.Item {
	items := make([]list.Item, len(dcs))
	for i, d := range dcs {
		items[i] = item{
			name:       d.Datacenter,
			totalVMs:   len(d.VMs),
			totalHosts: len(d.Hosts),
		}
	}

	return items
}
