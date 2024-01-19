package vm

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hsarena/vcbox/pkg/tui/common"
	"github.com/hsarena/vcbox/pkg/vmware"
)

func (bv BubbleVM) Update(msg tea.Msg) (BubbleVM, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		horizontal, vertical := common.ListStyle.GetFrameSize()
		paginatorHeight := lipgloss.Height(bv.list.Paginator.View())
		bv.list.SetSize(msg.Width-horizontal, msg.Height-vertical-paginatorHeight)
		bv.viewport = viewport.New(msg.Width, msg.Height)
		bv.viewport.SetContent(bv.detailView())
	}

	var cmd tea.Cmd
	bv.list, cmd = bv.list.Update(msg)

	return bv, cmd
}

func (bv BubbleVM) UpdateList(inventory []vmware.VMInventory, msg tea.Msg) (BubbleVM, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		horizontal, vertical := common.ListStyle.GetFrameSize()
		paginatorHeight := lipgloss.Height(bv.list.Paginator.View())
		bv.list.SetSize(msg.Width-horizontal, msg.Height-vertical-paginatorHeight)
		bv.viewport = viewport.New(msg.Width, msg.Height)
		bv.viewport.SetContent(bv.detailView())
	}
	return bv, bv.list.SetItems(vmToItem(inventory))
}

func (bv BubbleVM) MockUpdateList(inventory []vmware.MockVMInventory, msg tea.Msg) (BubbleVM, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		horizontal, vertical := common.ListStyle.GetFrameSize()
		paginatorHeight := lipgloss.Height(bv.list.Paginator.View())
		bv.list.SetSize(msg.Width-horizontal, msg.Height-vertical-paginatorHeight)
		bv.viewport = viewport.New(msg.Width, msg.Height)
		bv.viewport.SetContent(bv.detailView())
	}
	return bv, bv.list.SetItems(mockVMToItem(inventory))
}
