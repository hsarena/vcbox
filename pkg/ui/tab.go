package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	infos = iota
	logs
	events
	monitoring
	remote
	nLines int32 = 100
)

func NewTab() *tab {
	tab := &tab{
		self:       nil,
		infos:      nil,
		events:     nil,
		logs:       nil,
		monitoring: nil,
		remote:     nil,
	}
	tab.initTab()
	return tab
}

func (t *tab) initTab() {
	tab := tview.NewFlex()
	tab.SetBorder(false)
	infos := tview.NewTextView().SetText("")
	infos.SetBackgroundColor(tcell.Color(tcell.ColorDefault))
	events := tview.NewTextView().SetText("This is event tab")
	events.SetBackgroundColor(tcell.ColorDefault)
	logs := tview.NewTextView().SetText("This is log tab")
	logs.SetBackgroundColor(tcell.ColorDefault)
	monitoring := tview.NewFlex().AddItem(tview.NewTextView().SetText(""), 0, 1, false)
	monitoring.SetBackgroundColor(tcell.ColorDefault)
	remote := tview.NewFlex().AddItem(tview.NewTextView().SetText(""), 0, 1, false)
	remote.SetBackgroundColor(tcell.ColorDefault)
	t.self = tab
	t.infos = infos
	t.events = events
	t.logs = logs
	t.monitoring = monitoring
	t.remote = remote
}

func (t *tab) tabLayout(selected int) *tview.Flex {
	tab := tview.NewFlex().SetDirection(tview.FlexRow)
	tab.SetBorder(true)
	tab.SetBackgroundColor(tcell.ColorDefault)
	switch selected {
	case infos:
		tab.SetTitle(" Info ")
		tab.AddItem(t.infos, 0, 10, false)
	case logs:
		tab.SetTitle(" Logs ")
		tab.AddItem(t.logs, 0, 10, false)
	case events:
		tab.SetTitle(" Events ")
		tab.AddItem(t.events, 0, 10, false)
	case monitoring:
		tab.SetTitle(" Monitoring ")
		tab.AddItem(t.monitoring, 0, 10, false)
	case remote:
		tab.SetTitle(" Remote ")
		tab.AddItem(t.remote, 0, 10, false)
	}
	return tab
}

func (t *tab) Clear() {
	for i := 0; i < 5; i++ {
		switch i {
		case infos:
			t.infos.Clear()
		case logs:
			t.logs.Clear()
		case events:
			t.events.Clear()
		case monitoring:
			t.monitoring.Clear()
		case remote:
			t.remote.Clear()
		}
	}
}
