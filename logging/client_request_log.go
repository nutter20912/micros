package logging

import (
	"context"
	"encoding/json"

	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func PublishLog(ctx context.Context, p client.Message) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("topic", p.Topic()),
		attribute.String("ContentType", p.ContentType()))

	md, _ := metadata.FromContext(ctx)
	payload, _ := json.Marshal(map[string]interface{}{
		"metadata": md,
	})

	span.AddEvent("request",
		trace.WithAttributes(attribute.String("request.payload", string(payload))))
}
