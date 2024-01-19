package mock

type MockService struct{}

type Inventory struct {
	Datacenter string
	Hosts      []HostInventory
	VMs        []VMInventory
}

type HostInventory struct {
	Log             string
	ComputeResource string
}

type VMInventory struct {
	Name   string
	CPU    int32
	Memory int32
	OS     string
	IP     string
	Status string
}
