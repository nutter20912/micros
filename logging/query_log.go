package logging

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func MongoStartedLog(ctx context.Context, e *event.CommandStartedEvent) {
	span := trace.SpanFromContext(ctx)

	attrs := []attribute.KeyValue{
		semconv.DBSystemMongoDB,
		semconv.NetTransportTCP,
		semconv.DBOperation(e.CommandName),
		semconv.DBName(e.DatabaseName)}

	span.SetAttributes(attrs...)

	queryBytes, _ := bson.MarshalExtJSON(e.Command, false, false)
	span.SetAttributes(semconv.DBStatement(string(queryBytes)))
}

func MongoFinishedLog(ctx context.Context, e event.CommandFinishedEvent, err error) {
	span := trace.SpanFromContext(ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
}
