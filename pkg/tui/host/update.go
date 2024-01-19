package host

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
)

func (bh BubbleHost) Update(msg tea.Msg) (BubbleHost, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		horizontal, vertical := common.ListStyle.GetFrameSize()
		paginatorHeight := lipgloss.Height(bh.list.Paginator.View())
		bh.list.SetSize(msg.Width-horizontal, msg.Height-vertical-paginatorHeight)
		bh.viewport = viewport.New(msg.Width, msg.Height)
		bh.viewport.SetContent(bh.logView())
	}

	var cmd tea.Cmd
	bh.list, cmd = bh.list.Update(msg)

	return bh, cmd
}
