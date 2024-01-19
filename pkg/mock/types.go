package mock

import (
	"github.com/hsarena/vcbox/pkg/tui/datacenter"
	"github.com/hsarena/vcbox/pkg/tui/host"
	"github.com/hsarena/vcbox/pkg/tui/vm"
	"github.com/hsarena/vcbox/pkg/vmware"
)

type MockService struct{}

type state int
type model struct {
	bd            datacenter.BubbleDatacenter
	bh            host.BubbleHost
	bv            vm.BubbleVM
	inventory     []vmware.MockInventory
	selectedDC    int
	selectedHost  int
	selectedVM    int
	width, height int
	state         state
}
