package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/hsarena/vcbox/pkg/tui/side"
	"github.com/hsarena/vcbox/pkg/tui/tab"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"
)

type Model struct {
	header    viewport.Model
	side      side.Model
	tab       tab.Model
	inventory []vmware.Inventory
	selected  types.ManagedObjectReference
}

func InitModel(client *govmomi.Client) Model {
	d := vmware.NewDiscoveryService(client)
	m := vmware.NewMetricsService(client)
	inventory, err := d.FetchInventory()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	side := side.InitModel(inventory)
	tab := tab.InitModel(m)
	return Model{
		side:      side,
		tab:       tab,
		inventory: inventory,
	}
}
