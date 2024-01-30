package vmware

import (
	"context"
	"flag"
	"fmt"

	"github.com/hsarena/vcbox/pkg/util"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/types"
)

var interval = flag.Int("i", 300, "Interval ID")

func NewMetricsService(client *govmomi.Client) *MetricsService {
	return &MetricsService{client: client}
}

func (m *MetricsService) FetchMetrics(obj types.ManagedObjectReference, metrics []string) (mms map[string][]float64, err error) {
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
	counters, err := perfManager.CounterInfoByName(ctx)
	if err != nil {
		return nil, err
	}

	// Create PerfQuerySpec
	spec := types.PerfQuerySpec{
		Entity:     obj,
		MaxSample:  24,
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
	// Read result
	var (
		metricStr string
		name      string
		value     []float64
	)
	mms = make(map[string][]float64, len(result))
	for _, metric := range result {
		for _, v := range metric.Value {
			counter := counters[v.Name]
			units := counter.UnitInfo.GetElementDescription().Label

			if len(v.Value) != 0 {
				metricStr += fmt.Sprintf("\n ---| %s\t%s\t%s\n",
					v.Name, v.ValueCSV(), units)
				switch units {
				case "MHz":
					name = fmt.Sprintf("%s|%s", v.Name, "GHz")
					value = util.ToF64(v.Value, 1024)
				case "KB":
					name = fmt.Sprintf("%s|%s", v.Name, "GB")
					value = util.ToF64(v.Value, 1024*1024)
				case "%":
					name = fmt.Sprintf("%s|%s", v.Name, units)
					value = util.ToF64(v.Value, 100)
				case "num":
					name = fmt.Sprintf("%s|%s", v.Name, units)
					value = util.ToF64(v.Value, 1)
				case "KBps":
					name = fmt.Sprintf("%s|%s", v.Name, "MBps")
					value = util.ToF64(v.Value, 1024)
				}
				mms[name] = value
			}
		}
	}
	return mms, nil
}
