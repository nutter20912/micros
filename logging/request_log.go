package logging

import (
	"context"
	"encoding/json"

	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func RequestLog(ctx context.Context, req server.Request, rsp interface{}, err error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("service", req.Service()),
		attribute.String("method", req.Method()),
		attribute.String("endpoint", req.Endpoint()),
		attribute.String("contentType", req.ContentType()))

	md, _ := metadata.FromContext(ctx)
	payload, _ := json.Marshal(map[string]interface{}{
		"metadata": md,
		"header":   req.Header(),
		"body":     req.Body(),
	})

	span.AddEvent("request",
		trace.WithAttributes(attribute.String("request.payload", string(payload))))

	if err == nil {
		rspBytes, _ := json.Marshal(rsp)
		span.AddEvent("response",
			trace.WithAttributes(attribute.String("response.message", string(rspBytes))))
	}
}
