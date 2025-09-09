package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

// github.com/uptrace/uptrace-go
const meterName = "go.opentelemetry.io/contrib/examples/prometheus"

func main() {

	ctx := context.Background()

	// New returns a Prometheus Exporter.
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	// NewMeterProvider returns a new and configured MeterProvider.
	// By default, the returned MeterProvider is configured with the default Resource and no Readers.
	// Readers cannot be added after a MeterProvider is created.
	// This means the returned MeterProvider, one created with no Readers, will perform no operations.

	// WithReader associates Reader r with a MeterProvider.
	// By default, if this option is not used, the MeterProvider will perform no operations; no data will be exported without a Reader.
	meterProvider := metric.NewMeterProvider(metric.WithReader(exporter)).Meter(meterName)

	go server()

	measurementOption := otelMetric.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)

	float64Counter, err := meterProvider.Float64Counter("foo", otelMetric.WithDescription("foo description"))
	if err != nil {
		log.Fatal(err)
	}
	float64Counter.Add(ctx, 5, measurementOption)
}

func server() {
	log.Print("server address localhost:2333")
	http.Handle("/metters", promhttp.Handler())
	if err := http.ListenAndServe(":2333", nil); err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
