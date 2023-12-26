package box

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

type DCVMMap map[string][]*object.VirtualMachine

func DCDiscovery(c *govmomi.Client) (dcs []*object.Datacenter, err error) {
	// Create a view of Datastore objects
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create a new finder that will discover the defaults and are looked for Networks/Datastores
	f := find.NewFinder(c.Client, true)
	dcs, err = f.DatacenterList(ctx,"*")

	if err != nil {
		return nil, fmt.Errorf("couldn't discover any Datacenter instance inside of vCenter %v", err)
	}
	return dcs, nil
}

//VMInventory will create an inventory
func VMDiscoveryInsideDC(c *govmomi.Client, dcName string) ([]*object.VirtualMachine, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new finder that will discover the defaults and are looked for Networks/Datastores
	f := find.NewFinder(c.Client, true)

	// Find one and only datacenter, not sure how VMware linked mode will work
	dc, err := f.DatacenterOrDefault(ctx, dcName)
	if err != nil {
		return nil, fmt.Errorf("no Datacenter instance could be found inside of vCenter %v", err)
	}

	// Make future calls local to this datacenter
	f.SetDatacenter(dc)

	vms, err := f.VirtualMachineList(ctx, "*")

	if err != nil {
		return nil, fmt.Errorf("no VM could be found inside of datacenter %s %v", dcName, err)
	}
	// Sort function to sort by name
	sort.Slice(vms, func(i, j int) bool {
		switch strings.Compare(vms[i].Name(), vms[j].Name()) {
		case -1:
			return true
		case 1:
			return false
		}
		return vms[i].Name() > vms[j].Name()
	})

	return vms, nil
}


// DiscoverDCVMMap discovers datacenters and their corresponding VMs and returns a map
func DiscoverDCVMMap(c *govmomi.Client) (DCVMMap, error) {
	dcVMMap := make(DCVMMap)

	datacenters, err := DCDiscovery(c)
	if err != nil {
		return nil, err
	}

	for _, dc := range datacenters {
		vms, err := VMDiscoveryInsideDC(c, dc.Name())
		if err != nil {
			return nil, err
		}
		dcVMMap[dc.Name()] = vms
	}

	return dcVMMap, nil
}
