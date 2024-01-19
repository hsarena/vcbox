package tv

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewSide() *side {
	side := &side{
		datacenters: nil,
		hosts:       nil,
		vms:         nil,
	}
	side.initSide()
	return side
}

func (s *side) initSide() {
	dList := tview.NewList().ShowSecondaryText(false)
	dList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Datacenters ")
	dList.AddItem("", "", 0, nil)
	s.datacenters = dList

	hostList := tview.NewList().ShowSecondaryText(false)
	hostList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Compute Resources ")
	hostList.AddItem("", "", 0, nil)
	s.hosts = hostList

	vList := tview.NewList().ShowSecondaryText(false)
	vList.SetBorder(true).SetBackgroundColor(tcell.ColorDefault).
		SetTitle(" Virtual Machines ")
	vList.AddItem("", "", 0, nil)
	s.vms = vList

}

func (s *side) sideLayout() *tview.Flex {
	side := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(s.datacenters, 0, 1, true).
		AddItem(s.hosts, 0, 1, true).
		AddItem(s.vms, 0, 4, true)

	return side
}
