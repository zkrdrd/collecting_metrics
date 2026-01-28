package metrics

import (
	"runtime"
)

// MemoryMetrics структура для метрик памяти
type MemoryMetrics struct {
	TotalAlloc  uint64 `json:"total_alloc"`
	HeapAlloc   uint64 `json:"heap_alloc"`
	HeapSys     uint64 `json:"heap_sys"`
	StackInuse  uint64 `json:"stack_inuse"`
	StackSys    uint64 `json:"stack_sys"`
	MSpanInuse  uint64 `json:"mspan_inuse"`
	MSpanSys    uint64 `json:"mspan_sys"`
	MCacheInuse uint64 `json:"mcache_inuse"`
	MCacheSys   uint64 `json:"mcache_sys"`
}

// CPUMetrics структура для метрик ЦП
type CPUMetrics struct {
	Goroutines int     `json:"goroutines"`
	NumCPU     int     `json:"num_cpu"`
	Usage      float64 `json:"usage"`
}

// RuntimeMetrics структура для метрик рантайма
type RuntimeMetrics struct {
	NumGC         uint32  `json:"num_gc"`
	PauseTotalNs  uint64  `json:"pause_total_ns"`
	LastGC        uint64  `json:"last_gc"`
	NextGC        uint64  `json:"next_gc"`
	GCCPUFraction float64 `json:"gc_cpu_fraction"`
}

// CollectSystemMetrics собирает все системные метрики
func CollectSystemMetrics() (MemoryMetrics, CPUMetrics, RuntimeMetrics) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memMetrics := MemoryMetrics{
		TotalAlloc:  memStats.TotalAlloc,
		HeapAlloc:   memStats.HeapAlloc,
		HeapSys:     memStats.HeapSys,
		StackInuse:  memStats.StackInuse,
		StackSys:    memStats.StackSys,
		MSpanInuse:  memStats.MSpanInuse,
		MSpanSys:    memStats.MSpanSys,
		MCacheInuse: memStats.MCacheInuse,
		MCacheSys:   memStats.MCacheSys,
	}

	cpuMetrics := CPUMetrics{
		Goroutines: runtime.NumGoroutine(),
		NumCPU:     runtime.NumCPU(),
		Usage:      0.0, // Заполняется отдельно
	}

	runtimeMetrics := RuntimeMetrics{
		NumGC:         memStats.NumGC,
		PauseTotalNs:  memStats.PauseTotalNs,
		LastGC:        memStats.LastGC,
		NextGC:        memStats.NextGC,
		GCCPUFraction: memStats.GCCPUFraction,
	}

	return memMetrics, cpuMetrics, runtimeMetrics
}
