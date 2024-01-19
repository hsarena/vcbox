package host

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi/object"
)

type item struct {
	name string
	logs *object.DiagnosticLog
}

func (i item) Name() string        { return i.name }
func (i item) FilterValue() string { return i.name }

type BubbleHost struct {
	list     list.Model
	viewport viewport.Model
}

func NewBubbleHost(l list.Model) BubbleHost {
	return BubbleHost{list: l}
}

func InitialModel(inventory []vmware.HostInventory) BubbleHost {
	items := hostToItem(inventory)
	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = "Hosts"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	return BubbleHost{list: l}
}

func hostToItem(hosts []vmware.HostInventory) []list.Item {
	items := make([]list.Item, len(hosts))
	for i, h := range hosts {
		items[i] = item{
			name: h.ComputeResource.Name(),
			logs: h.Log,
		}
	}

	return items
}

func MockInitialModel(inventory []vmware.MockHostInventory) BubbleHost {
	items := mockHostToItem(inventory)
	l := list.New(items, itemDelegate{}, 0, 0)
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
