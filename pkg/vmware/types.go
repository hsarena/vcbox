package vmware

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
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

type VCClient struct {
	Address  string
	Username string
	Password string
	Insecure bool
}
