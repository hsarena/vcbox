package mock

import (
	"fmt"
)

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) FetchInventory() ([]Inventory, error) {
	dcD := [3]string{"HWB", "AST", "AFR"}
	in := make([]Inventory, len(dcD))
	for i, dc := range dcD {
		in[i].Datacenter = dc
		crD := [3]string{"srv-01", "srv-02", "srv-03"}
		hostI := make([]HostInventory, len(crD))
		for i, c := range crD {
			hostI[i].ComputeResource = c
			hostI[i].Log = fmt.Sprintf("This is the mock log of host [%s]", c)
		}
		vms := [10]string{"VM-01",
			"VM-02",
			"VM-03",
			"VM-04",
			"VM-05",
			"VM-06",
			"VM-07",
			"VM-08",
			"VM-09",
			"VM-10"}

		vmI := make([]VMInventory, len(vms))
		for i := range vms {
			vm := &VMInventory{
				Name:   vms[i],
				OS:     "openSUSE",
				CPU:    4,
				Memory: 4096,
				IP:     "192.168.10.100",
				Status: "On",
			}
			vmI[i] = *vm
		}
		in[i].Hosts = hostI
		in[i].VMs = vmI
	}
	return in, nil
}
