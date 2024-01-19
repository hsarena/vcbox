package tui

func (m model) View() string {
	switch m.state {

	case showDatacenterView:
		return m.bd.View(true)
	case showHostView:
		return m.bh.View(true)
	case showVMView:
		return m.bv.View(true)
	default:
		return m.bd.View(true)
	}
	//return m.bd.View(false) + m.bh.View(false) + m.bv.View(false)
}
