// Package tracing provides a tracing utility for the application
package tracing

import (
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("Docman"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

// SetupTracer sets the exporter and the resource for tracing
func SetupTracer() *trace.TracerProvider {
		// Write telemetry data to a file.
		f, err := os.Create("traces.txt")
		if err != nil {
			log.Error().Err(err)
		}
		// defer f.Close()
	
		exp, err := newExporter(f)
		if err != nil {
			log.Error().Err(err)
		}
	
		tp := trace.NewTracerProvider(
			trace.WithBatcher(exp),
			trace.WithResource(newResource()),
		)
		return tp

}