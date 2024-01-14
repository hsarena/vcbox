package ui

import (
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/rivo/tview"
)

type UI struct {
	app          *tview.Application
	page         *page
	side         *side
	selectedDc   int
	selectedHost int
	selectedVm   int
	inventory    []vmware.Inventory
}

type page struct {
	navbar *navbar
	tab    *tab
}

type tab struct {
	infos       *tview.TextView
	logs        *tview.TextView
	events      *tview.TextView
	monitoring  *tview.Flex
	remote      *tview.Flex
	selectedTab int
}

type navbar struct {
	bar *tview.Flex
}

type side struct {
	datacenters *tview.List
	hosts       *tview.List
	vms         *tview.List
}
