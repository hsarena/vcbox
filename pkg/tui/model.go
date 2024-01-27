package tui

import (
	"log"

	"github.com/hsarena/vcbox/pkg/tui/datacenter"
	"github.com/hsarena/vcbox/pkg/tui/host"
	"github.com/hsarena/vcbox/pkg/tui/vm"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/vmware/govmomi"
)

type model struct {
	bd            datacenter.BubbleDatacenter
	bh            host.BubbleHost
	bv            vm.BubbleVM
	inventory     []vmware.Inventory
	selectedDC    int
	selectedHost  int
	selectedVM    int
	width, height int
	state         state
}

func InitialModel(client *govmomi.Client) model {
	d := vmware.NewDiscoveryService(client)
	m := vmware.NewMetricsService(client)
	inventory, err := d.FetchInventory()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	return model{
		state:        showDatacenterView,
		bd:           datacenter.InitialModel(inventory),
		bh:           host.InitialModel(inventory[0].Hosts, m),
		bv:           vm.InitialModel(inventory[0].VMs, m),
		selectedDC:   0,
		selectedHost: 0,
		selectedVM:   0,
		inventory:    inventory,
	}
}
