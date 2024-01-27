package common

import "github.com/charmbracelet/lipgloss"

type ShowViewType int

const (
	ShowList ShowViewType = iota
	ShowDetail
	ShowFull
)

var (
	ListStyle = lipgloss.NewStyle().
			Width(30).
			PaddingRight(1).
			MarginRight(1).
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#f56642"))
	ListColorStyle = lipgloss.NewStyle().
			Background(lipgloss.NoColor{}).
			Foreground(lipgloss.Color("#00ffa2")).
			Underline(true)
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)
	ListSelectedListItemStyle = lipgloss.NewStyle().
					PaddingLeft(2).
					Foreground(lipgloss.Color("#00ffa2"))
	DetailStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(1)
	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f56642")).
			PaddingTop(1)
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffa2")).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#f56642")).
			Bold(true)

	HostMetrics = []string{"cpu.usagemhz.average",
		"mem.consumed.average"}
	VMMetrics = []string{"cpu.usage.average",
		"mem.usage.average",
		"net.usage.average",
		"datastore.maxTotalLatency.latest",
	}
)
