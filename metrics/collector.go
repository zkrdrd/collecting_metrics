package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"
)

const (
	TypeMetricFloat64Gauge = iota
	TypeMetricFloat64Counter
	TypeMetricInt64Gauge
	TypeMetricInt64Counter
)

type MetricsCollector struct {
	meter           metric.Meter
	cpuPercentUsage metric.Float64ObservableGauge
	memoryByteUsage metric.Int64ObservableGauge
	goroutineCount  metric.Int64ObservableGauge
	gcPauseTotal    metric.Float64ObservableCounter
	heapAlloc       metric.Int64ObservableGauge
	stackInUse      metric.Int64ObservableGauge
	gcCount         metric.Int64ObservableCounter
	mallocCount     metric.Int64ObservableCounter

	// Дополнительные поля
	// serviceName     string
	// environment     string
	// instanceID      string
	// startTime       time.Time
	// lastCPUTime     uint64
	// lastCPUCallTime time.Time
}

type TypeMetric int
type MetricUnit struct {
	Name        string
	Description string
	Instrument  string
	Type        TypeMetric
	Observer    metric.Observable
}

var cpuPercentUsage = &MetricUnit{
	Name:        `process.cpu.usage`,
	Description: `CPU usage percentage`,
	Instrument:  `%`,
	Type:        TypeMetricFloat64Gauge,
}

var memoryByteUsage = &MetricUnit{
	Name:        `process.memory.usage`,
	Description: `Memory usage in bytes`,
	Instrument:  `By`,
	Type:        TypeMetricInt64Gauge,
}

var goroutineCount = &MetricUnit{
	Name:        `process.runtime.go.goroutines`,
	Description: `Number of goroutines`,
	Instrument:  `1`,
	Type:        TypeMetricInt64Gauge,
}

var gcPauseTotal = &MetricUnit{
	Name:        `process.runtime.go.gc.pause.total`,
	Description: `Total GC pause duration`,
	Instrument:  `s`,
	Type:        TypeMetricFloat64Counter,
}

var heapAlloc = &MetricUnit{
	Name:        `process.runtime.go.mem.heap.alloc`,
	Description: `Bytes of allocated heap objects`,
	Instrument:  `By`,
	Type:        TypeMetricInt64Gauge,
}

var stackInUse = &MetricUnit{
	Name:        `process.runtime.go.mem.stack.inuse`,
	Description: `Bytes in stack spans`,
	Instrument:  `By`,
	Type:        TypeMetricInt64Gauge,
}

var gcCount = &MetricUnit{
	Name:        `process.runtime.go.gc.count`,
	Description: `Total number of GC cycles`,
	Instrument:  `1`,
	Type:        TypeMetricInt64Counter,
}

var mallocCount = &MetricUnit{
	Name:        `process.runtime.go.mem.malloc.count`,
	Description: `Cumulative count of heap objects allocated`,
	Instrument:  `1`,
	Type:        TypeMetricInt64Counter,
}

var ListMetrics = []*MetricUnit{
	cpuPercentUsage,
	memoryByteUsage,
	goroutineCount,
	gcPauseTotal,
	heapAlloc,
	stackInUse,
	gcCount,
	mallocCount,
}

type MetricUnits []*MetricUnit

// InitMetrics инициализирует все метрики
func (mu MetricUnits) PrepareMetrics(meter metric.Meter) error {
	var (
		o   metric.Observable
		err error
	)

	for _, unit := range mu {
		switch unit.Type {
		case TypeMetricFloat64Gauge:
			o, err = meter.Float64ObservableGauge(
				unit.Name,
				metric.WithDescription(unit.Description),
				metric.WithUnit(unit.Instrument),
			)
		case TypeMetricFloat64Counter:
			o, err = meter.Float64ObservableCounter(
				unit.Name,
				metric.WithDescription(unit.Description),
				metric.WithUnit(unit.Instrument),
			)
		case TypeMetricInt64Gauge:
			o, err = meter.Int64ObservableGauge(
				unit.Name,
				metric.WithDescription(unit.Description),
				metric.WithUnit(unit.Instrument),
			)
		case TypeMetricInt64Counter:
			o, err = meter.Int64ObservableCounter(
				unit.Name,
				metric.WithDescription(unit.Description),
				metric.WithUnit(unit.Instrument),
			)
		}

		if err != nil {
			return fmt.Errorf("prepare observer for unit %s, %s, %s: %w", unit.Name, unit.Type, unit.Instrument, err)
		}

		unit.Observer = o

	}

	return nil
}

func (mc MetricsCollector) InitMetricsUnits(units MetricUnits) error {
	if err := units.PrepareMetrics(mc.meter); err != nil {
		return fmt.Errorf(`init metrics: %w`, err)
	}
	return nil
}
