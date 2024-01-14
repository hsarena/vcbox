package internal

import "github.com/vmware/govmomi/object"

type dcInventory struct {
	datacenter *object.Datacenter
	hosts      []hostInventory
	vms        []vmInventory
}

type hostInventory struct {
	log             *object.DiagnosticLog
	computeResource *object.ComputeResource
}
type vmInventory struct {
	name   string
	cpu    int32
	memory int32
	os     string
	ip     string
	status string
}
