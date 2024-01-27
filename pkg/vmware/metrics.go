package vmware

import (
	"context"
	"flag"
	"log"

	// "fmt"

	"github.com/vmware/govmomi"
	// "github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/types"
)

var interval = flag.Int("i", 300, "Interval ID")

func NewMetricsService(client *govmomi.Client) *MetricsService {
	return &MetricsService{client: client}
}

func (m *MetricsService) FetchMetrics(obj types.ManagedObjectReference, metrics []string) (mapd map[string]performance.MetricSeries, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Get virtual machines references
	viewManager := view.NewManager(m.client.Client)

	v, err := viewManager.CreateContainerView(ctx, m.client.Client.ServiceContent.RootFolder, nil, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Create a PerfManager
	perfManager := performance.NewManager(m.client.Client)

	// Retrieve counters name list
	// counters, err := perfManager.CounterInfoByName(ctx)
	// if err != nil {
	// 	return nil,err
	// }
	// //var names []string
	// names := []string{
	// // "cpu.usage.average",
	// // "net.usage.average",
	// // "disk.write.average",
	// // "disk.read.average",
	// // "mem.usage.average",
	// "cpu.usagemhz.average",
	// "mem.overhead.average"}

	// // for name := range counters {
	// // 	names = append(names, name)
	// // }
	// Create PerfQuerySpec
	spec := types.PerfQuerySpec{
		Entity:     obj,
		MaxSample:  30,
		MetricId:   []types.PerfMetricId{{Instance: "*"}},
		IntervalId: int32(*interval),
	}

	// Query metrics
	sample, err := perfManager.SampleByName(ctx, spec, metrics, []types.ManagedObjectReference{obj})
	if err != nil {
		return nil, err
	}

	result, err := perfManager.ToMetricSeries(ctx, sample)
	if err != nil {
		return nil, err
	}
	log.Printf("the result metrics is: %v", result)
	// Read result
	mapd = make(map[string]performance.MetricSeries, len(result))
	for _, metric := range result {
		// // name := metric.Entity
		// vm := object.NewHostSystem(m.client.Client, metric.Entity)
		// name, err := vm.ObjectName(ctx)
		// if err != nil {
		// 	return nil, err
		// }
		for _, v := range metric.Value {
			// counter := counters[v.Name]
			// units := counter.UnitInfo.GetElementDescription().Label

			// instance := v.Instance
			// if instance == "" {
			// 	instance = "-"
			// }
			// log.Println("about to create mapd")
			if len(v.Value) != 0 {
				// str += fmt.Sprintf("%s\t%s\t%s\n",
				// 	v.Name, v.ValueCSV(), units)
				mapd[v.Name] = v
			}
		}
	}
	return mapd, nil
}
