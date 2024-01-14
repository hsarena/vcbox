package ui

import (
	"context"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/hsarena/vcbox/pkg/vmware"
	"github.com/rivo/tview"
	"github.com/vmware/govmomi"
)

const (
	info = iota
	logs
	events
	monitoring
	remote
	nLines int32 = 100
)

func NewUI(client *govmomi.Client) *UI {
	d := vmware.NewDiscoveryService(client)
	inventory, err := d.FetchInventory()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	app := tview.NewApplication()
	page := NewPage()
	side := NewSide()

	ui := &UI{
		app:          app,
		page:         page,
		side:         side,
		selectedDc:   0,
		selectedHost: 0,
		selectedVm:   0,
		inventory:    inventory,
	}

	ui.createUI()

	return ui
}

func (ui *UI) createUI() {

	ui.updateDatacenterList()
	ui.setupEventHandlers()
	ui.setupTabEventHandlers()
	ui.app.SetRoot(ui.mainLayout(), true).SetFocus(ui.side.datacenters).EnableMouse(false)
}

func (ui *UI) mainLayout() *tview.Flex {

	return tview.NewFlex().
		AddItem(ui.side.sideLayout(), 0, 1, true).
		AddItem(ui.page.tabLayout(), 0, 5, true)
}

func (ui *UI) updateDatacenterList() {

	ui.side.datacenters.Clear()

	for _, d := range ui.inventory {
		ui.side.datacenters.AddItem(d.Datacenter.Name(), "", 0, nil)
	}
}

func (ui *UI) updateComputeResourceList() {
	ui.side.hosts.Clear()

	for _, cr := range ui.inventory[ui.selectedDc].Hosts {
		ui.side.hosts.AddItem(cr.ComputeResource.Name(), "", 0, nil)
	}

}

func (ui *UI) updateVMList() {
	ui.side.vms.Clear()

	for _, vm := range ui.inventory[ui.selectedDc].VMs {
		ui.side.vms.AddItem(vm.Name, "", 0, nil)
	}

}

func (ui *UI) updateVMInfo() {
	ui.page.tab.infos.Clear()

	vmDetailsText := fmt.Sprintf("[orange]Name: [white]%s\n[orange]CPU: [white]%d\n[orange]Memory: [white]%d MB\n[orange]OS: [white]%s\n[orange]IP: [white]%s\n[orange]Status: [white]%s",
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].Name,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].CPU,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].Memory,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].OS,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].IP,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].Status)

	ui.page.tab.infos.SetText(vmDetailsText).SetDynamicColors(true)
}

func (ui *UI) updateHostLog() {
	ui.page.tab.logs.Clear()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := ui.inventory[ui.selectedDc].Hosts[ui.selectedHost].Log.Seek(ctx, nLines)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.inventory[ui.selectedDc].Hosts[ui.selectedHost].Log.Copy(ctx, ui.page.tab.logs)
}

func (ui *UI) setupEventHandlers() {

	// Datacenter Events
	ui.side.datacenters.SetFocusFunc(func() {
		ui.updateVMList()
		ui.updateComputeResourceList()
	})

	ui.side.datacenters.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.datacenters.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedDc = i
		ui.updateComputeResourceList()
		ui.updateVMList()
	})

	ui.side.datacenters.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.Stop()
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.side.vms)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.side.hosts)
		}
		return event
	})

	//Compute Resource Events
	ui.side.hosts.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.hosts.SetSelectedTextColor(tcell.ColorDarkOrange)
	})

	ui.side.hosts.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.hosts.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedHost = i
		ui.updateHostLog()
		ui.app.SetFocus(ui.page.tab.logs)
	})

	ui.side.hosts.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.side.datacenters)
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.side.hosts)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.side.vms)
		}
		return event
	})

	ui.side.vms.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.vms.SetSelectedTextColor(tcell.ColorDarkGreen)
	})
	ui.side.vms.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.vms.SetSelectedTextColor(tcell.ColorDarkGreen)
		ui.selectedVm = i
		ui.updateVMInfo()
		ui.app.SetFocus(ui.page.tab.infos)

	})

	ui.side.vms.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.side.datacenters)
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.side.vms)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.side.datacenters)
		}
		return event
	})
}

func (ui *UI) setupTabEventHandlers() {

	ui.page.tab.logs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.side.hosts)
		}
		return event
	})

}

func (ui *UI) Run() error {
	return ui.app.Run()
}
