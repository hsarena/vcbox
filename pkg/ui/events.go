package ui

import "github.com/gdamore/tcell/v2"

func (ui *UI) dcEventHandlers() {
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
}

func (ui *UI) hostEventHandlers() {
	ui.side.hosts.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.hosts.SetSelectedTextColor(tcell.ColorDarkOrange)
	})

	ui.side.hosts.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.hosts.SetSelectedTextColor(tcell.ColorDarkOrange)
		ui.selectedHost = i
		ui.selectedTab = logs
		ui.updateTab()
		ui.updateHostLog()
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
}

func (ui *UI) vmEventHandlers() {
	ui.side.vms.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.vms.SetSelectedTextColor(tcell.ColorDarkGreen)
	})
	ui.side.vms.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		ui.side.vms.SetSelectedTextColor(tcell.ColorDarkGreen)
		ui.selectedVm = i
		ui.selectedTab = infos
		ui.updateTab()
		ui.updateVMInfo()
	})
	ui.side.vms.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.side.datacenters)
		case tcell.KeyEnter:
			ui.app.SetFocus(ui.side.vms)
		case tcell.KeyTab:
			ui.app.SetFocus(ui.side.datacenters)
		case tcell.KeyCtrlR:
			ui.selectedTab = remote
			ui.remoteToHost()
		}
		return event
	})
}

func (ui *UI) setupEventHandlers() {
	ui.dcEventHandlers()
	ui.hostEventHandlers()
	ui.vmEventHandlers()
}

func (ui *UI) logsTabEventHandlers() {
	ui.tab.logs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.side.hosts)
		}
		return event
	})
}

func (ui *UI) infosTabEventHandlers() {
	ui.tab.infos.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.SetFocus(ui.side.vms)
		}
		return event
	})
}

func (ui *UI) setupTabEventHandlers() {
	ui.logsTabEventHandlers()
	ui.infosTabEventHandlers()
}
