package datacenter

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
)

func (bd BubbleDatacenter) Update(msg tea.Msg) (BubbleDatacenter, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		horizontal, vertical := common.ListStyle.GetFrameSize()
		paginatorHeight := lipgloss.Height(bd.list.Paginator.View())
		bd.list.SetSize(msg.Width-horizontal, msg.Height/8-vertical-paginatorHeight)
		bd.viewport = viewport.New(msg.Width, msg.Height)
		bd.viewport.SetContent(bd.detailView())
	}

	var cmd tea.Cmd
	bd.list, cmd = bd.list.Update(msg)
	return bd, cmd
}

func (bd *BubbleDatacenter) GetSelectedItem() int {
	return bd.list.Index()
}
