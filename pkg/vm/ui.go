package vm

import (
	"fmt"
	"log"

	//"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

type UI struct {
	app              *tview.Application
	datacenters      *tview.List
	computeResources *tview.List
	vms              *tview.List
	vmDetails        *tview.TextView
	discovery        *DiscoveryService
	selectedDc       *object.Datacenter
	selectedCr       *object.ComputeResource
	selectedVm       *object.VirtualMachine
}

func NewUI(client *govmomi.Client) *UI {
	app := tview.NewApplication()
	discovery := NewDiscoveryService(client)
	ui := &UI{
		app:              app,
		datacenters:      nil,
		computeResources: nil,
		vms:              nil,
		vmDetails:        nil,
		selectedDc:       nil,
		selectedCr:       nil,
		selectedVm:       nil,
		discovery:        discovery,
	}

	ui.createUI()

	return ui
}

func (ui *UI) createUI() {
	dList := tview.NewList().ShowSecondaryText(false)
	dList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Datacenters")
	dList.AddItem("", "", 0, nil)
	ui.datacenters = dList

	crList := tview.NewList().ShowSecondaryText(false)
	crList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Compute Resources ")
	crList.AddItem("", "", 0, nil)
	ui.computeResources = crList

	vList := tview.NewList().ShowSecondaryText(false)
	vList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Virtual Machines ")
	vList.AddItem("", "", 0, nil)
	ui.vms = vList

	vmdList := tview.NewTextView().SetTextAlign(tview.AlignLeft)
	vmdList.SetBackgroundColor(tcell.ColorDefault)
	vmdList.SetDynamicColors(true).
		SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" VM Details ")
	vmdList.SetText("")
	ui.vmDetails = vmdList

	ui.updateDatacentersList()
	ui.setupEventHandlers()
	ui.app.SetRoot(ui.flexLayout(), true).SetFocus(ui.datacenters).EnableMouse(true)
}

func (ui *UI) flexLayout() *tview.Flex {
	return tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.datacenters, 0, 1, true).
			AddItem(ui.computeResources, 0, 1, true).
			AddItem(ui.vms, 0, 4, true),
			0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.vmDetails, 0, 5, true),
			0, 5, true)
}

func (ui *UI) updateDatacentersList() {
	// If the list is nil, create a new one
	if ui.datacenters == nil {
		ui.createUI()
		return
	}

	ui.datacenters.Clear()

	datacenters, err := ui.discovery.DiscoverDatacenters()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	for _, dc := range datacenters {
		ui.datacenters.AddItem(dc.Name(), "", 0, nil)
	}
}

func (ui *UI) updateComputeResourcesList(dc *object.Datacenter) {
	ui.computeResources.Clear()

	if dc == nil {
		return
	}

	crs, err := ui.discovery.DiscoverComputeResource(dc)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	for _, cr := range crs {
		ui.computeResources.AddItem(cr.Name(), "", 0, nil)
	}

}

func (ui *UI) updateVMsList(dc *object.Datacenter) {
	ui.vms.Clear()

	if dc == nil {
		return
	}

	vms, err := ui.discovery.DiscoverVMsInsideDC(dc)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	for _, vm := range vms {
		ui.vms.AddItem(vm.Name(), "", 0, nil)
	}

	//ui.app.Draw()
}

func (ui *UI) updateVMInfo(vm *object.VirtualMachine) {
	ui.vmDetails.Clear()

	if vm == nil {
		return
	}

	vmInfo, err := ui.discovery.DiscoverVMInfo(vm)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	vmDetailsText := fmt.Sprintf("Name: %s\nCPU: %d\nMemory: %d MB\nOS: %s\nIPs: %s\nStatus: %s",
		vmInfo.Name, vmInfo.CPU, vmInfo.Memory, vmInfo.OS, vmInfo.IPs, vmInfo.Status)

	ui.vmDetails.SetText(vmDetailsText)
	//ui.app.Draw()
}

func (ui *UI) setupEventHandlers() {

	datacenters, err := ui.discovery.DiscoverDatacenters()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.datacenters.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		if i < len(datacenters) {
			ui.selectedDc = datacenters[i]
			ui.updateComputeResourcesList(ui.selectedDc)
			ui.updateVMsList(ui.selectedDc)
		}
	})
	crs, err := ui.discovery.DiscoverComputeResource(ui.selectedDc)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.computeResources.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		if i < len(crs) {
			ui.selectedCr = crs[i]
		}
	})

	vms, err := ui.discovery.DiscoverVMsInsideDC(ui.selectedDc)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.computeResources.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		if i < len(vms) {
			ui.selectedVm = vms[i]
			ui.updateVMInfo(ui.selectedVm)
		}
	})

	ui.datacenters.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.Stop()
		}
		return event
	})

	ui.computeResources.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.datacenters)
		}
		return event
	})

	ui.vms.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.computeResources)
		}
		return event
	})

	ui.vmDetails.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.vms)
		}
		return event
	})
}

func (ui *UI) Run() error {
	return ui.app.Run()
}
