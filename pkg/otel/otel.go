package otel

import (
	"context"
	"fmt"

	"github.com/kofj/ipi/pkg/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracerHTTP(endpoint string, headers map[string]string) *sdktrace.TracerProvider {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if endpoint == "" {
		endpoint = "oodev.yhylk.com" //without trailing slash
	}

	otlptracehttp.NewClient()

	otlpHTTPExporter, err := otlptracehttp.New(context.TODO(),
		// otlptracehttp.WithInsecure(), // use http & not https
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithURLPath("/api/default/traces"),
		otlptracehttp.WithHeaders(headers), // use this if you want to pass custom headers
	)

	if err != nil {
		fmt.Println("Error creating HTTP OTLP exporter: ", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String("ipi-server"),
		semconv.ServiceVersionKey.String(version.GitVersion),
		attribute.String("Environment", "k3s"),
		attribute.String("BuildDate", version.BuildDate),
		attribute.String("GitCommit", version.GitCommit),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpHTTPExporter),
	)
	otel.SetTracerProvider(tp)

	return tp
}
