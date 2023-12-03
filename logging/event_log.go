package logging

import (
	"context"
	"encoding/json"

	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func EventLog(ctx context.Context, msg server.Message) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("topic", msg.Topic()),
		attribute.String("ContentType", msg.ContentType()))

	md, _ := metadata.FromContext(ctx)
	payload, _ := json.Marshal(map[string]interface{}{
		"metadata": md,
		"header":   msg.Header(),
		"body":     msg.Payload(),
	})

	span.AddEvent("request",
		trace.WithAttributes(attribute.String("request.payload", string(payload))))
}
