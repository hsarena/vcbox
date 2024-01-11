package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
)

type VMInfo struct {
	Name   string
	Host   string
	CPU    int32
	Memory int32
	OS     string
	IP    string
	Status string
}

func Infos(vmo mo.VirtualMachine) *VMInfo {

	return &VMInfo{
		Name:   vmo.Config.Name,
		CPU:    vmo.Summary.Config.NumCpu,
		Memory: vmo.Summary.Config.MemorySizeMB,
		OS:     vmo.Guest.GuestFullName,
		IP:    vmo.Guest.IpAddress,
		Status: string(vmo.Summary.Runtime.PowerState),
	}
}

type DiscoveryService struct {
	client *govmomi.Client
}

func NewDiscoveryService(client *govmomi.Client) *DiscoveryService {
	return &DiscoveryService{client: client}
}

// DiscoverDatacenters retrieves a list of datacenters.
func (d *DiscoveryService) DiscoverDatacenters() ([]*object.Datacenter, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finder := find.NewFinder(d.client.Client, true)
	datacenters, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("failed to discover datacenters: %v", err)
	}

	return datacenters, nil
}

// DiscoverComputeResource retrieves a list of compute resources.
func (d *DiscoveryService) DiscoverComputeResource(dc *object.Datacenter) ([]*object.ComputeResource, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if dc == nil {
		return nil, fmt.Errorf("the datacenter is nil")
	}

	finder := find.NewFinder(d.client.Client, true)
	finder.SetDatacenter(dc)
	crs, err := finder.ComputeResourceList(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("failed to discover hosts: %v", err)
	}
	return crs, err
}

// DiscoverVMsInsideDC retrieves a list of VMs inside a datacenter.
func (d *DiscoveryService) DiscoverVMsInsideDC(dc *object.Datacenter) ([]*object.VirtualMachine, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finder := find.NewFinder(d.client.Client, true)
	finder.SetDatacenter(dc)
	vms, err := finder.VirtualMachineList(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("failed to discover VMs inside datacenter %s: %v", dc.Name(), err)
	}
	return vms, nil
}

// DiscoverVMsInsideHost retrieves a list of VMs inside a host.
func (d *DiscoveryService) DiscoverVMsInsideHost(dc *object.Datacenter, h *object.HostSystem) ([]*object.VirtualMachine, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finder := find.NewFinder(d.client.Client, true)
	finder.SetDatacenter(dc)
	vms, err := finder.VirtualMachineList(ctx, "*")
	if err != nil {
		return nil, fmt.Errorf("failed to discover VMs inside datacenter %s: %v", dc.Name(), err)
	}
	return vms, nil
}

// DiscoverVMInfo retrieves details of a VM.
func (d *DiscoveryService) DiscoverVMInfo(vm *object.VirtualMachine) (*VMInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var vmMo mo.VirtualMachine

	err := object.NewVirtualMachine(d.client.Client, vm.Reference()).Properties(ctx, vm.Reference(), []string{"summary", "guest", "config", "runtime"}, &vmMo)
	if err != nil {
		return nil, fmt.Errorf("failed to discover VM info for %s: %v", vm.Name(), err)
	}

	if err != nil {
		log.Printf("failed to discover VM Host info for %s: %v", vm.Name(), err)
	}

	return Infos(vmMo), nil
}

