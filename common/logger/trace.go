package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceIDFromOTel(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return ""
	}
	return sc.TraceID().String()
}
