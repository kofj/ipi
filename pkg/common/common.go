package common

import "go.opentelemetry.io/otel"

var Tracer = otel.Tracer("github.com/kofj/ipi")
