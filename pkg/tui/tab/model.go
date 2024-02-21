package tab

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/hsarena/vcbox/pkg/vmware"
)

type tabKind int

const (
	Metrics tabKind = iota
	Logs
	Remote
)

var (
	tabKinds = []tabKind{Metrics, Logs, Remote}
)

type Model struct {
	view      viewport.Model
	tabs      []tab
	activeTab int
	metrics   *vmware.MetricsService
}

type tab struct {
	content string
	kind    tabKind
}

func newTab(c string, k tabKind) tab {
	return tab{
		content: c,
		kind:    k,
	}
}

func InitModel(metrics *vmware.MetricsService) Model {
	tabs := make([]tab, len(tabKinds))
	var c string
	for i, k := range tabKinds {
		switch k {
		case Metrics:
			c = "This is a tab for showing the metrics"
		case Logs:
			c = "This is a tab for showing the logs"
		case Remote:
			c = "This is a tab for showing remote console"
		}

		tabs[i] = newTab(c, k)

	}
	m := Model{
		tabs:    tabs,
		metrics: metrics,
	}
	return m
}
