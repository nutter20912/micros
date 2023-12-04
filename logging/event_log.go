package logging

import (
	"context"
	"encoding/json"

	"github.com/spf13/viper"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func EventLog(ctx context.Context, msg server.Message) {
	span := trace.SpanFromContext(ctx)

	brokerName := viper.GetString("micro.broker")
	attrs := []attribute.KeyValue{
		semconv.PeerServiceKey.String(brokerName),
		semconv.MessagingSystem(brokerName),
		semconv.MessagingSourceName(msg.Topic())}

	span.SetAttributes(attrs...)

	md, _ := metadata.FromContext(ctx)
	message, _ := json.Marshal(map[string]interface{}{
		"metadata":    md,
		"header":      msg.Header(),
		"contentType": msg.ContentType(),
		"body":        msg.Payload(),
	})

	span.AddEvent("request", trace.WithAttributes(
		attribute.String("message", string(message))))
}
