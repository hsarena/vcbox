package tui

import (
	"log"

	"github.com/hsarena/vcbox/pkg/tui/common"
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
	tabs          []common.Tab
	width, height int
	state         state
	activeTab     int
}

func InitialModel(client *govmomi.Client) model {
	d := vmware.NewDiscoveryService(client)
	m := vmware.NewMetricsService(client)
	inventory, err := d.FetchInventory()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	tab := make([]common.Tab, 2)
	tab[0].Name = "details"
	tab[0].Content = ""
	tab[1].Name = "metrics"
	tab[1].Content = ""
	return model{
		state:     showDatacenterView,
		bd:        datacenter.InitialModel(inventory, m),
		bh:        host.InitialModel(inventory[0].Hosts, m),
		bv:        vm.InitialModel(inventory[0].VMs, m),
		inventory: inventory,
		tabs:      tab,
	}
}
