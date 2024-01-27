package host

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

type item struct {
	name string
	logs *object.DiagnosticLog
	obj  types.ManagedObjectReference
}

func (i item) Name() string        { return i.name }
func (i item) FilterValue() string { return i.name }

type BubbleHost struct {
	list     list.Model
	viewport viewport.Model
	metrics  *vmware.MetricsService
}

func NewBubbleHost(l list.Model) BubbleHost {
	return BubbleHost{list: l}
}

func InitialModel(inventory []vmware.HostInventory, metrics *vmware.MetricsService) BubbleHost {
	items := hostToItem(inventory)
	l := list.New(items, hostItemDelegate{}, 0, 0)
	l.Title = "Hosts"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	return BubbleHost{list: l, metrics: metrics}
}

func hostToItem(hosts []vmware.HostInventory) []list.Item {
	items := make([]list.Item, len(hosts))
	for i, h := range hosts {
		items[i] = item{
			name: h.HostSystem.Name(),
			obj:  h.HostSystem.Reference(),
			logs: h.Log,
		}
	}

	return items
}

func MockInitialModel(inventory []vmware.MockHostInventory) BubbleHost {
	items := mockHostToItem(inventory)
	l := list.New(items, hostItemDelegate{}, 0, 0)
	l.Title = "Hosts"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	return BubbleHost{list: l}
}

func mockHostToItem(hosts []vmware.MockHostInventory) []list.Item {
	items := make([]list.Item, len(hosts))
	for i, h := range hosts {
		items[i] = item{
			name: h.ComputeResource,
		}
	}

	return items
}
