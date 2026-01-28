package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"collecting_metrics/metrics"
)

func main() {
	// Инициализация контекста
	// ctx := context.Background()

	// Инициализация Prometheus экспортера
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	}

	// Создание провайдера метрик
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)

	// Установка глобального провайдера
	otel.SetMeterProvider(provider)

	// Создание сборщика метрик
	collector := metrics.NewMetricsCollector(
		"golang-service-name",
		"production",
		"instance-1",
	)

	// Инициализация метрик
	if err := collector.InitMetrics(); err != nil {
		log.Fatalf("Failed to init metrics: %v", err)
	}

	// Регистрация callback'ов
	if err := collector.Register(); err != nil {
		log.Fatalf("Failed to register metrics: %v", err)
	}

	log.Println("Metrics collector initialized successfully")
	log.Printf("Service uptime: %v", collector.GetUptime())

	// Запуск HTTP сервера для Prometheus
	go func() {
		// Используем стандартный HTTP handler от Prometheus экспортера
		log.Println("Starting metrics server on :8080")

		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Демонстрационная нагрузка для сбора метрик
	go generateLoad()

	// Основной цикл приложения
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Uptime: %v", collector.GetUptime())

			// Сбор и вывод текущих метрик
			memMetrics, cpuMetrics, runtimeMetrics := metrics.CollectSystemMetrics()
			log.Printf("Goroutines: %d, HeapAlloc: %d MB",
				cpuMetrics.Goroutines,
				memMetrics.HeapAlloc/(1024*1024))
			log.Printf("GC Cycles: %d, Total Pause: %.2f ms",
				runtimeMetrics.NumGC,
				float64(runtimeMetrics.PauseTotalNs)/1e6)
		}
	}
}

// generateLoad создает нагрузку для демонстрации метрик
func generateLoad() {
	for i := 0; i < 100; i++ {
		go func(id int) {
			data := make([]byte, 1024*1024) // 1MB
			for {
				// Простая нагрузка
				for j := range data {
					data[j] = byte((j + id) % 256)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// Периодически создаем новые горутины
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	counter := 100
	for range ticker.C {
		go func(id int) {
			time.Sleep(30 * time.Second)
		}(counter)
		counter++
		if counter%10 == 0 {
			log.Printf("Created %d goroutines", counter)
		}
	}
}
