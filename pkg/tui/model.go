package tui

import (
	"github.com/hsarena/vcbox/pkg/tui/side"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi"
	"log"
)

type Model struct {
	side      side.Model
	inventory []vmware.Inventory
}

func InitModel(client *govmomi.Client) Model {
	d := vmware.NewDiscoveryService(client)
	inventory, err := d.FetchInventory()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	side := side.InitModel(inventory)
	return Model{
		side:      side,
		inventory: inventory,
	}
}
