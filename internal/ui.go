package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

type UI struct {
	app         *tview.Application
	dcsList     *tview.List
	crsList     *tview.List
	vmsList     *tview.List
	tabPage     *tab
	selectedDc  int
	selectedCr  int
	selectedVm  int
	dcInventory []dcInventory
}

type dcInventory struct {
	datacenter       *object.Datacenter
	computeResources []*object.ComputeResource
	virtualMachines  []vmInventory
	hostLog          *object.DiagnosticLog
}

type vmInventory struct {
	name   string
	cpu    int32
	memory int32
	os     string
	ip     string
	status string
}

type tab struct {
	navbar     *tview.Flex
	infos      *tview.TextView
	logs       *tview.TextView
	events     *tview.TextView
	monitoring *tview.Flex
	remote     *tview.Flex
}

func NewUI(client *govmomi.Client) *UI {
	app := tview.NewApplication()

	ui := &UI{
		app:        app,
		dcsList:    nil,
		crsList:    nil,
		vmsList:    nil,
		selectedDc: 0,
		selectedCr: 0,
		selectedVm: 0,
		tabPage:    nil,
	}

	ui.initMap(client)
	ui.createUI()

	return ui
}

func (ui *UI) initMap(client *govmomi.Client) {
	d := NewDiscoveryService(client)
	dcs, err := d.DiscoverDatacenters()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	dcI := make([]dcInventory, len(dcs))
	for i, dc := range dcs {
		dcI[i].datacenter = dc
		crs, err := d.DiscoverComputeResource(dc)
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		dcI[i].computeResources = crs
		dcI[i].hostLog, err = d.FetchHostLogs(dcI[i].computeResources[0])
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		vms, err := d.DiscoverVMsInsideDC(dc)
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		vm := make([]vmInventory, len(vms))
		for j, v := range vms {
			vmInfo, err := d.DiscoverVMInfo(v)
			if err != nil {
				log.Printf("%s", err.Error())
				return
			}
			vm[j].name = vmInfo.Name
			vm[j].cpu = vmInfo.CPU
			vm[j].memory = vmInfo.Memory
			vm[j].os = vmInfo.OS
			vm[j].ip = vmInfo.IP
			vm[j].status = vmInfo.Status
		}
		dcI[i].virtualMachines = vm
	}
	ui.dcInventory = dcI
	ui.selectedDc = 0
	ui.selectedCr = 0
}

func (ui *UI) initSide() {
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

}

func (ui *UI) initTab() {

	infos := tview.NewTextView().SetText("")
	infos.SetBackgroundColor(tcell.ColorDefault)
	infos.SetBorder(true)
	events := tview.NewTextView().SetText("This is event page")
	events.SetBackgroundColor(tcell.ColorDefault)
	logs := tview.NewTextView().SetText("This is log page")
	logs.SetBackgroundColor(tcell.ColorDefault)
	monitoring := tview.NewFlex().AddItem(tview.NewTextView().SetText(""), 0, 1, false)
	monitoring.SetBackgroundColor(tcell.ColorDefault)
	remote := tview.NewFlex().AddItem(tview.NewTextView().SetText(""), 0, 1, false)
	remote.SetBackgroundColor(tcell.ColorDefault)
	navbar := tview.NewFlex().SetDirection(tview.FlexColumn)
	navbar.SetBackgroundColor(tcell.ColorDarkBlue)
	iText := tview.NewTextView().SetDynamicColors(true)
	iText.SetBorder(true)
	iText.SetBackgroundColor(tcell.ColorDefault)
	iText.SetText("[orange](I)[green]nfo")
	lText := tview.NewTextView().SetDynamicColors(true)
	lText.SetBorder(true)
	lText.SetBackgroundColor(tcell.ColorDefault)
	lText.SetText("[orange](L)[green]ogs")
	eText := tview.NewTextView().SetDynamicColors(true)
	eText.SetBorder(true)
	eText.SetBackgroundColor(tcell.ColorDefault)
	eText.SetText("[orange](E)[green]vents")
	mText := tview.NewTextView().SetDynamicColors(true)
	mText.SetBorder(true)
	mText.SetBackgroundColor(tcell.ColorDefault)
	mText.SetText("[orange](M)[green]onitoring")
	rText := tview.NewTextView().SetDynamicColors(true)
	rText.SetBorder(true)
	rText.SetBackgroundColor(tcell.ColorDefault)
	rText.SetText("[orange](R)[green]emote")
	navbar.AddItem(iText, 0, 1, false).
		AddItem(lText, 0, 1, false).
		AddItem(eText, 0, 1, false).
		AddItem(mText, 0, 1, false).
		AddItem(rText, 0, 1, false)

	ui.tabPage = &tab{
		navbar:     navbar,
		infos:      infos,
		events:     events,
		logs:       logs,
		monitoring: monitoring,
		remote:     remote,
	}

}

func (ui *UI) createUI() {
	ui.initTab()
	ui.initSide()
	ui.updateDatacenterList()
	ui.setupEventHandlers()
	ui.app.SetRoot(ui.flexLayout(), true).SetFocus(ui.dcsList).EnableMouse(false)
}

func (ui *UI) flexLayout() *tview.Flex {
	tabPage := tview.NewFlex().SetDirection(tview.FlexRow)
	tabPage.SetBorder(true)
	tabPage.SetBackgroundColor(tcell.ColorDefault)
	tabPage.AddItem(ui.tabPage.navbar, 0, 1, false).
		AddItem(ui.tabPage.logs, 0, 10, false)

	sidePage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.dcsList, 0, 1, true).
		AddItem(ui.crsList, 0, 1, true).
		AddItem(ui.vmsList, 0, 4, true)
	return tview.NewFlex().
		AddItem(sidePage, 0, 1, true).
		AddItem(tabPage, 0, 5, true)
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

func (ui *UI) updateVMsList() {
	ui.vmsList.Clear()

	for _, vm := range ui.dcInventory[ui.selectedDc].virtualMachines {
		ui.vmsList.AddItem(vm.name, "", 0, nil)
	}

}

func (ui *UI) updateVMInfo() {
	ui.tabPage.infos.Clear()

	vmDetailsText := fmt.Sprintf("[orange]Name: [white]%s\n[orange]CPU: [white]%d\n[orange]Memory: [white]%d MB\n[orange]OS: [white]%s\n[orange]IP: [white]%s\n[orange]Status: [white]%s",
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].name,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].cpu,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].memory,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].os,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].ip,
		ui.dcInventory[ui.selectedDc].virtualMachines[ui.selectedVm].status)

	ui.tabPage.infos.SetText(vmDetailsText).SetDynamicColors(true)
}

func (ui *UI) updateHostLog() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := ui.dcInventory[ui.selectedDc].hostLog.Seek(ctx, nLines)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.dcInventory[ui.selectedDc].hostLog.Copy(ctx, ui.tabPage.logs)
}

func (ui *UI) setupEventHandlers() {

	// Datacenter Events
	ui.dcsList.SetFocusFunc(func() {
		ui.updateVMsList()
		ui.updateComputeResourceList()
	})

	ui.dcsList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.dcsList.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedDc = i
		ui.updateComputeResourceList()
		ui.updateVMsList()
	})

	ui.dcsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.Stop()
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.vmsList)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.crsList)
		}
		return event
	})

	//Compute Resource Events
	ui.crsList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.crsList.SetSelectedTextColor(tcell.ColorDarkOrange)
	})

	ui.crsList.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.crsList.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedCr = i
		log.Println("about to update host log")
		ui.updateHostLog()
		ui.app.SetFocus(ui.tabPage.logs)
	})

	ui.crsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.dcsList)
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.crsList)
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
		ui.app.SetFocus(ui.tabPage.infos)

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
