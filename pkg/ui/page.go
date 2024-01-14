package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewPage() *page {
	page := &page{
		navbar: nil,
		tab:    nil,
	}
	page.initNavbar()
	page.initTab()
	return page
}

func (p *page) initNavbar() {
	bar := tview.NewFlex().SetDirection(tview.FlexColumn)
	bar.SetBackgroundColor(tcell.ColorDarkBlue)
	for i := 0; i < 5; i++ {
		text := tview.NewTextView().SetDynamicColors(true)
		text.SetBorder(true)
		text.SetBackgroundColor(tcell.ColorDefault)
		switch i {
		case info:
			text.SetText("[orange](I)[green]nfo")
		case logs:
			text.SetText("[orange](L)[green]ogs")
		case events:
			text.SetText("[orange](E)[green]vents")
		case monitoring:
			text.SetText("[orange](M)[green]onitoring")
		case remote:
			text.SetText("[orange](R)[green]emote")
		}
		bar.AddItem(text, 0, 1, false)
	}
	p.navbar = &navbar{bar: bar}
}

func (p *page) initTab() {
	infos := tview.NewTextView().SetText("")
	infos.SetBackgroundColor(tcell.Color(tcell.ColorDefault))
	infos.SetBorder(true)
	events := tview.NewTextView().SetText("This is event page")
	events.SetBackgroundColor(tcell.ColorDefault)
	logs := tview.NewTextView().SetText("This is log page")
	logs.SetBackgroundColor(tcell.ColorDefault)
	monitoring := tview.NewFlex().AddItem(tview.NewTextView().SetText(""), 0, 1, false)
	monitoring.SetBackgroundColor(tcell.ColorDefault)
	remote := tview.NewFlex().AddItem(tview.NewTextView().SetText(""), 0, 1, false)
	remote.SetBackgroundColor(tcell.ColorDefault)

	p.tab = &tab{
		infos:       infos,
		events:      events,
		logs:        logs,
		monitoring:  monitoring,
		remote:      remote,
		selectedTab: 0, //infos
	}
}

func (p *page) tabLayout() *tview.Flex {
	tab := tview.NewFlex().SetDirection(tview.FlexRow)
	tab.SetBorder(true)
	tab.SetBackgroundColor(tcell.ColorDefault)
	tab.AddItem(p.navbar.bar, 0, 1, false)
	switch p.tab.selectedTab {
	case info:
		tab.AddItem(p.tab.infos, 0, 10, false)
	case logs:
		tab.AddItem(p.tab.logs, 0, 10, false)
	case events:
		tab.AddItem(p.tab.events, 0, 10, false)
	case monitoring:
		tab.AddItem(p.tab.monitoring, 0, 10, false)
	case remote:
		tab.AddItem(p.tab.remote, 0, 10, false)
	}
	return tab
}
