package metrics

import (
	"context"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// MetricsCollector структура для сбора метрик
type MetricsCollector struct {
	meter          metric.Meter
	cpuUsage       metric.Float64ObservableGauge
	memoryUsage    metric.Int64ObservableGauge
	goroutineCount metric.Int64ObservableGauge
	gcPauseTotal   metric.Float64ObservableCounter
	heapAlloc      metric.Int64ObservableGauge
	stackInUse     metric.Int64ObservableGauge
	gcCount        metric.Int64ObservableCounter
	mallocCount    metric.Int64ObservableCounter

	// Дополнительные поля
	serviceName     string
	environment     string
	instanceID      string
	startTime       time.Time
	lastCPUTime     uint64
	lastCPUCallTime time.Time
}

// NewMetricsCollector создает новый сборщик метрик
func NewMetricsCollector(serviceName, environment, instanceID string) *MetricsCollector {
	meter := otel.GetMeterProvider().Meter(
		"app.metrics",
		metric.WithInstrumentationVersion("1.0.0"),
	)

	return &MetricsCollector{
		meter:       meter,
		serviceName: serviceName,
		environment: environment,
		instanceID:  instanceID,
		startTime:   time.Now(),
	}
}

// InitMetrics инициализирует все метрики
func (mc *MetricsCollector) InitMetrics() error {
	var err error

	// CPU Usage
	mc.cpuUsage, err = mc.meter.Float64ObservableGauge(
		"process.cpu.usage",
		metric.WithDescription("CPU usage percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return err
	}

	// Memory Usage
	mc.memoryUsage, err = mc.meter.Int64ObservableGauge(
		"process.memory.usage",
		metric.WithDescription("Memory usage in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	// Goroutine Count
	mc.goroutineCount, err = mc.meter.Int64ObservableGauge(
		"process.runtime.go.goroutines",
		metric.WithDescription("Number of goroutines"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	// GC Pause Total
	mc.gcPauseTotal, err = mc.meter.Float64ObservableCounter(
		"process.runtime.go.gc.pause.total",
		metric.WithDescription("Total GC pause duration"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	// Heap Allocation
	mc.heapAlloc, err = mc.meter.Int64ObservableGauge(
		"process.runtime.go.mem.heap.alloc",
		metric.WithDescription("Bytes of allocated heap objects"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	// Stack In Use
	mc.stackInUse, err = mc.meter.Int64ObservableGauge(
		"process.runtime.go.mem.stack.inuse",
		metric.WithDescription("Bytes in stack spans"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	// GC Count
	mc.gcCount, err = mc.meter.Int64ObservableCounter(
		"process.runtime.go.gc.count",
		metric.WithDescription("Total number of GC cycles"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	// Malloc Count
	mc.mallocCount, err = mc.meter.Int64ObservableCounter(
		"process.runtime.go.mem.malloc.count",
		metric.WithDescription("Cumulative count of heap objects allocated"),
		metric.WithUnit("1"),
	)

	return err
}

// getCPUUsageSimplified упрощенный расчет использования CPU
func (mc *MetricsCollector) getCPUUsageSimplified() float64 {
	// Простой расчет на основе времени работы горутин
	// В реальном приложении лучше использовать github.com/shirou/gopsutil/cpu
	return float64(runtime.NumGoroutine()) / float64(runtime.NumCPU()) * 10.0
}

// CollectMetricsCallback функция обратного вызова для сбора метрик
func (mc *MetricsCollector) CollectMetricsCallback(ctx context.Context, o metric.Observer) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Сбор метрик ЦП
	cpuPercent := mc.getCPUUsageSimplified()

	// Атрибуты для всех метрик
	attrs := []attribute.KeyValue{
		attribute.String("service.name", mc.serviceName),
		attribute.String("environment", mc.environment),
		attribute.String("instance.id", mc.instanceID),
		attribute.String("go.version", runtime.Version()),
	}

	// CPU Usage
	o.ObserveFloat64(mc.cpuUsage, cpuPercent, metric.WithAttributes(attrs...))

	// Memory Usage (общее использование)
	o.ObserveInt64(mc.memoryUsage, int64(m.Sys), metric.WithAttributes(attrs...))

	// Goroutine Count
	o.ObserveInt64(mc.goroutineCount, int64(runtime.NumGoroutine()),
		metric.WithAttributes(attrs...))

	// GC Pause Total (наносекунды -> секунды)
	o.ObserveFloat64(mc.gcPauseTotal, float64(m.PauseTotalNs)/1e9,
		metric.WithAttributes(attrs...))

	// Heap Allocation
	o.ObserveInt64(mc.heapAlloc, int64(m.HeapAlloc),
		metric.WithAttributes(attrs...))

	// Stack In Use
	o.ObserveInt64(mc.stackInUse, int64(m.StackInuse),
		metric.WithAttributes(attrs...))

	// GC Count
	o.ObserveInt64(mc.gcCount, int64(m.NumGC),
		metric.WithAttributes(attrs...))

	// Malloc Count
	o.ObserveInt64(mc.mallocCount, int64(m.Mallocs),
		metric.WithAttributes(attrs...))

	return nil
}

// Register регистрирует сборщик метрик
func (mc *MetricsCollector) Register() error {
	_, err := mc.meter.RegisterCallback(
		mc.CollectMetricsCallback,
		mc.cpuUsage,
		mc.memoryUsage,
		mc.goroutineCount,
		mc.gcPauseTotal,
		mc.heapAlloc,
		mc.stackInUse,
		mc.gcCount,
		mc.mallocCount,
	)
	return err
}

// GetUptime возвращает время работы приложения
func (mc *MetricsCollector) GetUptime() time.Duration {
	return time.Since(mc.startTime)
}
