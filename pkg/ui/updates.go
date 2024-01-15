package ui

import (
	"context"
	"fmt"
	"log"
)

func (ui *UI) updateTab() {
	ui.tab.self.Clear()
	ui.tab.self.AddItem(ui.tab.tabLayout(ui.selectedTab), 0, 1, true)
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
	ui.tab.infos.Clear()
	vmDetailsText := fmt.Sprintf("[orange]Name: [white]%s\n[orange]CPU: [white]%d\n[orange]Memory: [white]%d MB\n[orange]OS: [white]%s\n[orange]IP: [white]%s\n[orange]Status: [white]%s",
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].Name,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].CPU,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].Memory,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].OS,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].IP,
		ui.inventory[ui.selectedDc].VMs[ui.selectedVm].Status)

	ui.tab.infos.SetText(vmDetailsText).SetDynamicColors(true)
}

func (ui *UI) updateHostLog() {
	ui.tab.logs.Clear()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := ui.inventory[ui.selectedDc].Hosts[ui.selectedHost].Log.Seek(ctx, nLines)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	ui.inventory[ui.selectedDc].Hosts[ui.selectedHost].Log.Copy(ctx, ui.tab.logs)
}
