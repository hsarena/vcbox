package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vmware/govmomi"
)

type UI struct {
	app            *tview.Application
	datacenterList *tview.List
	hostList       *tview.List
	vmList         *tview.List
	tabPage        *tab
	selectedDc     int
	selectedHost   int
	selectedVm     int
	dcInventory    []dcInventory
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
		app:            app,
		datacenterList: nil,
		hostList:       nil,
		vmList:         nil,
		selectedDc:     0,
		selectedHost:   0,
		selectedVm:     0,
		tabPage:        nil,
	}

	ui.initMap(client)
	ui.createUI()

	return ui
}

func (ui *UI) initMap(client *govmomi.Client) {
	d := NewDiscoveryService(client)
	dcD, err := d.DiscoverDatacenters()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	dcI := make([]dcInventory, len(dcD))
	for i, dc := range dcD {
		dcI[i].datacenter = dc
		crD, err := d.DiscoverComputeResource(dc)
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		hostI := make([]hostInventory, len(crD))
		for i, c := range crD {
			hostI[i].computeResource = c
			hostI[i].log, err = d.FetchHostLogs(c)
			if err != nil {
				log.Printf("%s", err.Error())
				return
			}
		}
		vms, err := d.DiscoverVMsInsideDC(dc)
		if err != nil {
			log.Printf("%s", err.Error())
			return
		}
		vmI := make([]vmInventory, len(vms))
		for i, v := range vms {
			vmInfo, err := d.DiscoverVMInfo(v)
			if err != nil {
				log.Printf("%s", err.Error())
				return
			}
			vmI[i].name = vmInfo.Name
			vmI[i].cpu = vmInfo.CPU
			vmI[i].memory = vmInfo.Memory
			vmI[i].os = vmInfo.OS
			vmI[i].ip = vmInfo.IP
			vmI[i].status = vmInfo.Status
		}
		dcI[i].hosts = hostI
		dcI[i].vms = vmI
	}
	ui.dcInventory = dcI
	ui.selectedDc = 0
	ui.selectedHost = 0
}

func (ui *UI) initSide() {
	dList := tview.NewList().ShowSecondaryText(false)
	dList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Datacenters ")
	dList.AddItem("", "", 0, nil)
	ui.datacenterList = dList

	hostList := tview.NewList().ShowSecondaryText(false)
	hostList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Compute Resources ")
	hostList.AddItem("", "", 0, nil)
	ui.hostList = hostList

	vList := tview.NewList().ShowSecondaryText(false)
	vList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Virtual Machines ")
	vList.AddItem("", "", 0, nil)
	ui.vmList = vList

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
	ui.setupTabEventHandlers()
	ui.app.SetRoot(ui.flexLayout(), true).SetFocus(ui.datacenterList).EnableMouse(false)
}

func (ui *UI) flexLayout() *tview.Flex {
	tabPage := tview.NewFlex().SetDirection(tview.FlexRow)
	tabPage.SetBorder(true)
	tabPage.SetBackgroundColor(tcell.ColorDefault)
	tabPage.AddItem(ui.tabPage.navbar, 0, 1, false).
		AddItem(ui.tabPage.logs, 0, 10, false)

	sidePage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.datacenterList, 0, 1, true).
		AddItem(ui.hostList, 0, 1, true).
		AddItem(ui.vmList, 0, 4, true)
	return tview.NewFlex().
		AddItem(sidePage, 0, 1, true).
		AddItem(tabPage, 0, 5, true)
}

func (ui *UI) updateDatacenterList() {

	ui.datacenterList.Clear()

	for _, d := range ui.dcInventory {
		ui.datacenterList.AddItem(d.datacenter.Name(), "", 0, nil)
	}
}

func (ui *UI) updateComputeResourceList() {
	ui.hostList.Clear()

	for _, cr := range ui.dcInventory[ui.selectedDc].hosts {
		ui.hostList.AddItem(cr.computeResource.Name(), "", 0, nil)
	}

}

func (ui *UI) updateVMList() {
	ui.vmList.Clear()

	for _, vm := range ui.dcInventory[ui.selectedDc].vms {
		ui.vmList.AddItem(vm.name, "", 0, nil)
	}

}

func (ui *UI) updateVMInfo() {
	ui.tabPage.infos.Clear()

	vmDetailsText := fmt.Sprintf("[orange]Name: [white]%s\n[orange]CPU: [white]%d\n[orange]Memory: [white]%d MB\n[orange]OS: [white]%s\n[orange]IP: [white]%s\n[orange]Status: [white]%s",
		ui.dcInventory[ui.selectedDc].vms[ui.selectedVm].name,
		ui.dcInventory[ui.selectedDc].vms[ui.selectedVm].cpu,
		ui.dcInventory[ui.selectedDc].vms[ui.selectedVm].memory,
		ui.dcInventory[ui.selectedDc].vms[ui.selectedVm].os,
		ui.dcInventory[ui.selectedDc].vms[ui.selectedVm].ip,
		ui.dcInventory[ui.selectedDc].vms[ui.selectedVm].status)

	ui.tabPage.infos.SetText(vmDetailsText).SetDynamicColors(true)
}

func (ui *UI) updateHostLog() {
	ui.tabPage.logs.Clear()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := ui.dcInventory[ui.selectedDc].hosts[ui.selectedHost].log.Seek(ctx, nLines)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.dcInventory[ui.selectedDc].hosts[ui.selectedHost].log.Copy(ctx, ui.tabPage.logs)
}

func (ui *UI) setupEventHandlers() {

	// Datacenter Events
	ui.datacenterList.SetFocusFunc(func() {
		ui.updateVMList()
		ui.updateComputeResourceList()
	})

	ui.datacenterList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.datacenterList.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedDc = i
		ui.updateComputeResourceList()
		ui.updateVMList()
	})

	ui.datacenterList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.Stop()
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.vmList)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.hostList)
		}
		return event
	})

	//Compute Resource Events
	ui.hostList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.hostList.SetSelectedTextColor(tcell.ColorDarkOrange)
	})

	ui.hostList.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.hostList.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedHost = i
		log.Println("about to update host log")
		ui.updateHostLog()
		ui.app.SetFocus(ui.tabPage.logs)
	})

	ui.hostList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.datacenterList)
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.hostList)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.vmList)
		}
		return event
	})

	ui.vmList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.vmList.SetSelectedTextColor(tcell.ColorDarkGreen)
	})
	ui.vmList.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.vmList.SetSelectedTextColor(tcell.ColorDarkGreen)
		ui.selectedVm = i
		ui.updateVMInfo()
		ui.app.SetFocus(ui.tabPage.infos)

	})

	ui.vmList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.datacenterList)
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.vmList)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.datacenterList)
		}
		return event
	})
}

func (ui *UI) setupTabEventHandlers() {

	ui.tabPage.logs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.hostList)
		}
		return event
	})

}

func (ui *UI) Run() error {
	return ui.app.Run()
}
