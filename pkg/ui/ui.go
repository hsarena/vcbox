package ui

import (
	"log"

	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/rivo/tview"
	"github.com/vmware/govmomi"
)

func NewUI(client *govmomi.Client) *UI {
	d := vmware.NewDiscoveryService(client)
	inventory, err := d.FetchInventory()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	app := tview.NewApplication()
	tab := NewTab()
	side := NewSide()

	ui := &UI{
		app:          app,
		tab:          tab,
		side:         side,
		selectedDc:   0,
		selectedHost: 0,
		selectedVm:   0,
		selectedTab:  infos,
		inventory:    inventory,
	}
	ui.createUI()
	return ui
}

func (ui *UI) createUI() {
	ui.updateDatacenterList()
	ui.setupEventHandlers()
	ui.setupTabEventHandlers()
	main := tview.NewFlex().
		AddItem(ui.side.sideLayout(), 0, 1, true).
		AddItem(ui.tab.self, 0, 5, true)
	ui.app.SetRoot(main, true).SetFocus(ui.side.datacenters).EnableMouse(false)
}

func (ui *UI) Run() error {
	return ui.app.Run()
}
