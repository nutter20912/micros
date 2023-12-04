package logging

import (
	"context"
	"encoding/json"

	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func PublishLog(ctx context.Context, c client.Client, p client.Message) {
	span := trace.SpanFromContext(ctx)

	broker := c.Options().Broker
	attrs := []attribute.KeyValue{
		semconv.PeerServiceKey.String(broker.String()),
		semconv.MessagingSystem(broker.String()),
		semconv.MessagingDestinationName(p.Topic()),

		attribute.String("contentType", c.Options().ContentType)}

	span.SetAttributes(attrs...)

	md, _ := metadata.FromContext(ctx)
	message, _ := json.Marshal(map[string]interface{}{
		"metadata":    md,
		"topic":       p.Topic(),
		"contentType": p.ContentType(),
		"payload":     p.Payload(),
	})

	span.AddEvent("request", trace.WithAttributes(
		attribute.String("message", string(message))))
}
