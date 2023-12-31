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
	app          *tview.Application
	datacenters  *tview.List
	vms          *tview.List
	vmDetails    *tview.TextView
	discovery    *DiscoveryService
	selectedDC   *object.Datacenter
	selectedVM   *object.VirtualMachine
}

func NewUI(client *govmomi.Client) *UI {
	app := tview.NewApplication()
	discovery := NewDiscoveryService(client)
	ui := &UI{
		app:         app,
		datacenters: nil,
		vms:         nil,
		vmDetails:   nil,
		discovery:   discovery,
		selectedDC:  nil,
		selectedVM:  nil,
	}

	ui.createUI()

	return ui
}

func (ui *UI) createUI() {
	log.Println("about to create ui")
	dList := tview.NewList().ShowSecondaryText(false)
	dList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle("Datacenters")
	dList.AddItem("", "", 0, nil)
	ui.datacenters = dList

	
	vList := tview.NewList().ShowSecondaryText(false)
	vList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle("Virtual Machines")
	vList.AddItem("", "", 0, nil)
	ui.vms = vList
	
	vmdList := tview.NewTextView().SetTextAlign(tview.AlignLeft)
	vmdList.SetDynamicColors(true).
		SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle("VM Details")
	vmdList.SetText("")
	ui.vmDetails = vmdList

	ui.updateDatacentersList()
	ui.setupEventHandlers()
	ui.app.SetRoot(ui.flexLayout(), true).SetFocus(ui.datacenters).EnableMouse(true)
}

func (ui *UI) flexLayout() *tview.Flex {
	return tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.datacenters, 0, 2, true).
			AddItem(ui.vms, 0, 4, true),
		0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.vmDetails,0,5,true),
		0, 5, true)
}

func (ui *UI) updateDatacentersList() {
	log.Println("about to update datacenter lists")
	// If the list is nil, create a new one
	if ui.datacenters == nil {
		log.Println("datacenter was empty then create ui")
		ui.createUI()
		return
	}


	log.Println("about to clear datacenter list")
	ui.datacenters.Clear()

	datacenters, err := ui.discovery.DiscoverDatacenters()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	for _, dc := range datacenters {
		log.Println("list of datacenters", dc.Name())
		ui.datacenters.AddItem(dc.Name(), "", 0, func() {
			ui.selectedDC = dc
			log.Println("selected datacenter is:", ui.selectedDC)
			log.Println("about to update vm lists")
			ui.updateVMsList()
			log.Println("about to update vm infos")
			ui.updateVMInfo()
		})
	}

	log.Println("about to draw app")
	//ui.app.Draw()
}

func (ui *UI) updateVMsList() {
	log.Println("about to update vms list")
	ui.vms.Clear()

	if ui.selectedDC == nil {
		return
	}

	vms, err := ui.discovery.DiscoverVMsInsideDC(ui.selectedDC)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	for _, vm := range vms {
		ui.vms.AddItem(vm.Name(), "", 0, func() {
			ui.selectedVM = vm
			ui.updateVMInfo()
		})
	}

	//ui.app.Draw()
}

func (ui *UI) updateVMInfo() {
	log.Println("about to update vm details")
	ui.vmDetails.Clear()

	if ui.selectedVM == nil {
		return
	}

	vmInfo, err := ui.discovery.DiscoverVMInfo(ui.selectedVM)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	vmDetailsText := fmt.Sprintf("Name: %s\nCPU: %d\nMemory: %d MB\nOS: %s\nIPs: %s\nStatus: %s",
		vmInfo.Name, vmInfo.CPU, vmInfo.Memory, vmInfo.OS, vmInfo.IPs, vmInfo.Status)

	ui.vmDetails.SetText(vmDetailsText)
	ui.app.Draw()
}

func (ui *UI) setupEventHandlers() {
	
	log.Println("about to setup event handler")
	ui.datacenters.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.Stop()
		}
		return event
	})

	datacenters, err := ui.discovery.DiscoverDatacenters()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.datacenters.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		if i < len(datacenters) {
			ui.selectedDC = datacenters[i]
			ui.updateVMsList()
		}
	})

	ui.vms.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.datacenters)
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
	log.Println("about to run ui")
	return ui.app.Run()
}
