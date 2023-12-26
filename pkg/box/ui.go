package box

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vmware/govmomi/object"
)
 
func NewUi(dcvm DCVMMap) {
	app := tview.NewApplication()
	d, v := VMBox(dcvm)
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(d,0,1,true).
			AddItem(v,0,6,true),
			0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetBackgroundColor(tcell.ColorDefault).SetTitle("Infos"),0,5,false),
			0, 5, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func DatacenterBox(dcs []*object.Datacenter) *tview.List {
	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).SetTitle("Datacenters")
	for _, dc := range dcs {
		list.AddItem(dc.Name(), "",0, nil)
	}
	return list
}

func VirtualMachinBox(vms []*object.VirtualMachine, dcName string) *tview.List {
	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).SetTitle("Virtual Machines")
	for _, vm := range vms {
		list.AddItem(vm.Name(), "",0, nil)
	}
	return list
}

func VMBox(dcvm DCVMMap) (*tview.List, *tview.List) {
	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).SetTitle("Datacenters")
	vmlist := tview.NewList().ShowSecondaryText(false)
	vmlist.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).SetTitle("Virtual Machines")
	for dc, vms := range dcvm {
		list.AddItem(dc, "",0, nil)
		for _, vm := range vms {

			vmlist.AddItem(vm.Name(), "",0, nil)
		}
	}
	return list, vmlist
}