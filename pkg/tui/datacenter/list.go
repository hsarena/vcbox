package datacenter

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hsarena/vcbox/pkg/tui/common"
)

type dcItemDelegate struct{}

func (d dcItemDelegate) Height() int                               { return 1 }
func (d dcItemDelegate) Spacing() int                              { return 0 }
func (d dcItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d dcItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	line := i.name

	if index == m.Index() {
		line = common.ListSelectedListItemStyle.Render("_|" + line)
	} else {
		line = common.ListItemStyle.Render(line)
	}

	fmt.Fprint(w, line)
}
