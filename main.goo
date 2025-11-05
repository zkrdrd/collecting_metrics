package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

// github.com/uptrace/uptrace-go
const meterName = "go.opentelemetry.io/contrib/examples/prometheus"

func main() {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

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

	// Meter returns a Meter with the given name and configured with options.
	// The name should be the name of the instrumentation scope creating telemetry.
	// This name may be the same as the instrumented code only if that code provides built-in instrumentation.
	// Calls to the Meter method after Shutdown has been called will return Meters that perform no operations.
	// This method is safe to call concurrently.
	meterProvider := metric.NewMeterProvider(metric.WithReader(exporter)).Meter(meterName)

	go server()

	// WithAttributes converts attributes into an attribute Set and sets the Set to be associated with a measurement.
	measurementOption := otelMetric.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)

	// Float64Counter is an instrument that records increasing float64 values.
	// Warning: Methods may be added to this interface in minor releases.
	// See package documentation on API implementation for information on how to set default behavior for unimplemented methods.

	// WithDescription sets the instrument description.
	float64Counter, err := meterProvider.Float64Counter("foo", otelMetric.WithDescription("foo description"))
	if err != nil {
		log.Fatal(err)
	}
	// Add records a change to the counter.
	// Use the WithAttributeSet (or, if performance is not a concern,
	// the WithAttributes) option to include measurement attributes.
	float64Counter.Add(ctx, 5, measurementOption)

	// Float64ObservableGauge returns a new Float64ObservableGauge instrument
	// identified by name and configured with options. The instrument is used
	// to asynchronously record instantaneous float64 measurements once per a
	// measurement collection cycle.
	//
	// Measurements for the returned instrument are made via a callback. Use
	// the WithFloat64Callback option to register the callback here, or use the
	// RegisterCallback method of this Meter to register one later. See the
	// Measurements section of the package documentation for more information.
	//
	// The name needs to conform to the OpenTelemetry instrument name syntax.
	// See the Instrument Name section of the package documentation for more
	// information.
	float64ObservableGauge, err := meterProvider.Float64ObservableGauge("bar", otelMetric.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}

	// RegisterCallback registers f to be called during the collection of a
	// measurement cycle.
	//
	// If Unregister of the returned Registration is called, f needs to be
	// unregistered and not called during collection.
	//
	// The instruments f is registered with are the only instruments that f may
	// observe values for.
	//
	// If no instruments are passed, f should not be registered nor called
	// during collection.
	//
	// The function f needs to be concurrent safe.
	_, err = meterProvider.RegisterCallback(func(_ context.Context, o otelMetric.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		// ObserveFloat64 records the float64 value for obsrv.
		o.ObserveFloat64(float64ObservableGauge, n, measurementOption)
		return nil
	}, float64ObservableGauge)
	if err != nil {
		log.Fatal(err)
	}

	// This is the equivalent of prometheus.NewHistogramVec
	// Float64Histogram returns a new Float64Histogram instrument identified by
	// name and configured with options. The instrument is used to
	// synchronously record the distribution of float64 measurements during a
	// computational operation.
	//
	// The name needs to conform to the OpenTelemetry instrument name syntax.
	// See the Instrument Name section of the package documentation for more
	// information.
	histogram, err := meterProvider.Float64Histogram(
		"baz",
		otelMetric.WithDescription("a histogram with custom buckets and rename"),

		//WithExplicitBucketBoundaries sets the instrument explicit bucket boundaries.
		//This option is considered "advisory", and may be ignored by API implementations.
		otelMetric.WithExplicitBucketBoundaries(64, 128, 256, 512, 1024, 2048, 4096),
	)
	if err != nil {
		log.Fatal(err)
	}

	histogram.Record(ctx, 136, measurementOption)
	histogram.Record(ctx, 64, measurementOption)
	histogram.Record(ctx, 701, measurementOption)
	histogram.Record(ctx, 830, measurementOption)

	// NotifyContext returns a copy of the parent context that is marked done (its Done channel is closed) when
	// one of the listed signals arrives, when the returned stop function is called,
	// or when the parent context's Done channel is closed, whichever happens first.
	// The stop function unregisters the signal behavior, which, like signal.
	// Reset, may restore the default behavior for a given signal.
	// For example, the default behavior of a Go program receiving os.Interrupt is to exit.
	// Calling NotifyContext(parent, os.Interrupt) will change the behavior to cancel the returned context.
	// Future interrupts received will not trigger the default (exit) behavior until the returned stop function is called.
	// The stop function releases resources associated with it,
	// so code should call stop as soon as the operations running in this Context complete and
	// signals no longer need to be diverted to the context.
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func server() {
	log.Print("server address localhost:2333")
	http.Handle("/metters", promhttp.Handler())
	if err := http.ListenAndServe(":2333", nil); err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
