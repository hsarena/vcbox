package vmware

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
)

type Inventory struct {
	Datacenter *object.Datacenter
	Hosts      []HostInventory
	VMs        []VMInventory
}

type HostInventory struct {
	Log             *object.DiagnosticLog
	ComputeResource *object.ComputeResource
}

type VMInventory struct {
	Name   string
	CPU    int32
	Memory int32
	OS     string
	IP     string
	Status string
}

type DiscoveryService struct {
	client *govmomi.Client
}

func NewVMInventory(vmo mo.VirtualMachine) *VMInventory {

	return &VMInventory{
		Name:   vmo.Config.Name,
		CPU:    vmo.Summary.Config.NumCpu,
		Memory: vmo.Summary.Config.MemorySizeMB,
		OS:     vmo.Guest.GuestFullName,
		IP:     vmo.Guest.IpAddress,
		Status: string(vmo.Summary.Runtime.PowerState),
	}
}

func NewDiscoveryService(client *govmomi.Client) *DiscoveryService {
	return &DiscoveryService{client: client}
}
