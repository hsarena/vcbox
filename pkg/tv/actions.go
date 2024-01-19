package tv

import (
	"github.com/hsarena/vcbox/pkg/ssh"
)

func (ui *UI) remoteToHost() {
	ui.tab.remote.Clear()
	ssh.SSHConnect("jadmin", ui.inventory[ui.selectedDc].VMs[ui.selectedVm].IP, 22, true, ui.tab.remote)
}
