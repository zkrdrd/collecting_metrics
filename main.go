package main

import (
	"log"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// github.com/uptrace/uptrace-go
const meterName = "go.opentelemetry.io/contrib/examples/prometheus"

func main() {

	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	meter := metric.NewMeterProvider(metric.WithReader(exporter)).Meter(meterName)
}
