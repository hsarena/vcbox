package vmware

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

type Inventory struct {
	Datacenter *object.Datacenter
	Hosts      []HostInventory
	VMs        []VMInventory
}

type HostInventory struct {
	CpuModel string
	NumCpuCores int16
	MemorySize int64
	Uptime int32
	NumNics int32
	NumHBAs int32
	HostMaxVirtualDiskCapacity int64
	PowerState string
	Log        *object.DiagnosticLog
	HostSystem *object.HostSystem
}

type VMInventory struct {
	Name   string
	CPU    int32
	Memory int32
	OS     string
	IP     string
	Status string
	VM     types.ManagedObjectReference
}

type DiscoveryService struct {
	client *govmomi.Client
}

type MetricsService struct {
	client *govmomi.Client
}
type VCClient struct {
	Address  string
	Username string
	Password string
	Insecure bool
}

type MockInventory struct {
	Datacenter string
	Hosts      []MockHostInventory
	VMs        []MockVMInventory
}

type MockHostInventory struct {
	Log             string
	ComputeResource string
}

type MockVMInventory struct {
	Name   string
	CPU    int32
	Memory int32
	OS     string
	IP     string
	Status string
}
