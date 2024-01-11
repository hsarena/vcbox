package internal

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

type UI struct {
	app              *tview.Application
	dcsList      *tview.List
	crsList *tview.List
	vmsList              *tview.List
	vmDetails        *tview.TextView
	selectedDc      int
	selectedCr       int
	selectedVm       int
	dcInventory		[]dcInventory
}

type dcInventory struct {
	datacenter *object.Datacenter
	computeResources []*object.ComputeResource
	virtualMachines []vmInventory
}

type vmInventory struct {
	name string
	cpu int32
	memory int32
	os string
	ip string
	status string
}

func NewUI(client *govmomi.Client) *UI {
	app := tview.NewApplication()
	
	ui := &UI{
		app:              app,
		dcsList:      nil,
		crsList: nil,
		vmsList:              nil,
		vmDetails:        nil,
		selectedDc:       0,
		selectedCr:       0,
		selectedVm:       0,
	}

	ui.initMap(client)
	ui.createUI()

	return ui
}

func (ui *UI) initMap(client *govmomi.Client) {
	d := NewDiscoveryService(client)
	dcs, err := d.DiscoverDatacenters()
	ui.dcInventory = make([]dcInventory, len(dcs))
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	for i, dc := range dcs {
		ui.dcInventory[i].datacenter = dc
		crs , err := d.DiscoverComputeResource(dc)
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		ui.dcInventory[i].computeResources = crs
		vms, err := d.DiscoverVMsInsideDC(dc)
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		ui.dcInventory[i].virtualMachines = make([]vmInventory, len(vms))
		for j, v := range vms {
			vmInfo, err := d.DiscoverVMInfo(v)
			if err != nil {
				log.Printf("%s", err.Error())
				return
			}
			ui.dcInventory[i].virtualMachines[j].name = vmInfo.Name
			ui.dcInventory[i].virtualMachines[j].cpu = vmInfo.CPU
			ui.dcInventory[i].virtualMachines[j].memory = vmInfo.Memory
			ui.dcInventory[i].virtualMachines[j].os = vmInfo.OS
			ui.dcInventory[i].virtualMachines[j].ip = vmInfo.IP
			ui.dcInventory[i].virtualMachines[j].status = vmInfo.Status
		}
	}
	ui.selectedDc = 0
}

func (ui *UI) createUI() {
	dList := tview.NewList().ShowSecondaryText(false)
	dList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Datacenters ")
	dList.AddItem("", "", 0, nil)
	ui.dcsList = dList

	crsList := tview.NewList().ShowSecondaryText(false)
	crsList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Compute Resources ")
	crsList.AddItem("", "", 0, nil)
	ui.crsList = crsList

	vList := tview.NewList().ShowSecondaryText(false)
	vList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Virtual Machines ")
	vList.AddItem("", "", 0, nil)
	ui.vmsList = vList

	vmdList := tview.NewTextView().SetTextAlign(tview.AlignLeft)
	vmdList.SetBackgroundColor(tcell.ColorDefault)
	vmdList.SetDynamicColors(true).
		SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" VM Details ")
	vmdList.SetText("")
	ui.vmDetails = vmdList

	ui.updateDatacenterList()
	ui.setupEventHandlers()
	ui.app.SetRoot(ui.flexLayout(), true).SetFocus(ui.dcsList).EnableMouse(false)
}

func (ui *UI) flexLayout() *tview.Flex {
	return tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.dcsList, 0, 1, true).
			AddItem(ui.crsList, 0, 1, false).
			AddItem(ui.vmsList, 0, 4, true),
			0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.vmDetails, 0, 5, false),
			0, 5, true)
}

func (ui *UI) updateDatacenterList() {

	ui.dcsList.Clear()

	for _, d := range ui.dcInventory {
		ui.dcsList.AddItem(d.datacenter.Name(), "", 0, nil)
	}
}

func (ui *UI) updateComputeResourceList() {
	ui.crsList.Clear()

	for _, cr := range ui.dcInventory[ui.selectedDc].computeResources {
		ui.crsList.AddItem(cr.Name(), "", 0, nil)
	}

}

func (ui *UI) updatevmsListList() {
	ui.vmsList.Clear()

	for _, vm := range ui.dcInventory[ui.selectedDc].virtualMachines {
		ui.vmsList.AddItem(vm.name, "", 0, nil)
	}

}


func (ui *UI) updateVMInfo() {
	ui.vmDetails.Clear()

	
	vmDetailsText := fmt.Sprintf("Name: %s\nCPU: %d\nMemory: %d MB\nOS: %s\nIPs: %s\nStatus: %s",
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].name,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].cpu,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].memory,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].os,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].ip,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].status)

	ui.vmDetails.SetText(vmDetailsText)
	//ui.app.Draw()
}


func (ui *UI) setupEventHandlers() {

	// Datacenter Events
	ui.dcsList.SetFocusFunc(func() {
		ui.updatevmsListList()
		ui.updateComputeResourceList()
	})

	ui.dcsList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.dcsList.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedDc = i
		ui.updateComputeResourceList()
		ui.updatevmsListList()
	})

	ui.dcsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.Stop()
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.vmsList)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.vmsList)
		}
		return event
	})

	ui.vmsList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.vmsList.SetSelectedTextColor(tcell.ColorDarkGreen)
	})
	ui.vmsList.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.vmsList.SetSelectedTextColor(tcell.ColorDarkGreen)
		ui.selectedVm = i
		ui.updateVMInfo()

	})

	ui.vmsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.dcsList)
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.vmsList)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.dcsList)
		}
		return event
	})

}

func (ui *UI) Run() error {
	return ui.app.Run()
}
