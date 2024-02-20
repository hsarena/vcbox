package side

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/hsarena/vcbox/pkg/vmware"
)

type Model struct {
	inventory  []vmware.Inventory
	sides      []side
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
			l = list.New(dcToItem(inventory), itemDelegate{}, 0, 7)
		case Host:
			l = list.New(hostToItem(inventory[0].Hosts), itemDelegate{}, 0, 8)
		case VM:
			l = list.New(vmToItem(inventory[0].VMs), itemDelegate{}, 0, 18)
		}

		sides[i] = newSide(l, k)

	}
	return Model{
		inventory:  inventory,
		sides:      sides,
		activeSide: Datacenter,
	}
}

func (m Model) GetSelectedItem() (list.Item, kind) {
	return m.sides[m.activeSide].list.SelectedItem(), m.activeSide
}
