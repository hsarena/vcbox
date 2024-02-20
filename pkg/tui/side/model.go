package side

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/hsarena/vcbox/pkg/vmware"
)

type Model struct {
	inventory  []vmware.Inventory
	sides      []side
	view       viewport.Model
	activeSide kind
}

type side struct {
	list list.Model
	kind kind
}

func newSide(l list.Model, k kind) side {
	return side{
		list: l,
		kind: k,
	}
}

func InitModel(inventory []vmware.Inventory) Model {
	sides := make([]side, len(kinds))
	var l list.Model
	for i, k := range kinds {
		switch k {
		case Datacenter:
			l = list.New(dcToItem(inventory), itemDelegate{}, 0, 0)
		case Host:
			l = list.New(hostToItem(inventory[0].Hosts), itemDelegate{}, 0, 0)
		case VM:
			l = list.New(vmToItem(inventory[0].VMs), itemDelegate{}, 0, 0)
		}

		sides[i] = newSide(l, k)

	}
	return Model{
		inventory:  inventory,
		sides:      sides,
		view: viewport.New(35,25),
		activeSide: Datacenter,
	}
}
