package host

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hsarena/vcbox/pkg/tui/common"
)

type hostItemDelegate struct{}

func (d hostItemDelegate) Height() int                               { return 1 }
func (d hostItemDelegate) Spacing() int                              { return 0 }
func (d hostItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d hostItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
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
