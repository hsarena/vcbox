package tv

import (
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/rivo/tview"
)

type UI struct {
	app          *tview.Application
	tab          *tab
	side         *side
	selectedDc   int
	selectedHost int
	selectedVm   int
	selectedTab  int
	inventory    []vmware.Inventory
}

type tab struct {
	self       *tview.Flex
	infos      *tview.TextView
	logs       *tview.TextView
	events     *tview.TextView
	monitoring *tview.TextView
	remote     *tview.TextView
}

type side struct {
	datacenters *tview.List
	hosts       *tview.List
	vms         *tview.List
}
