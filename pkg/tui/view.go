package tui

func (m model) View() string {
	switch m.state {

	case showDatacenterView:
		return m.bd.View()
	case showHostView:
		return m.bh.View()
	case showVMView:
		return m.bv.View()
	default:
		return m.bd.View()
	}
}
