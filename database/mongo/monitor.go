package mongo

import (
	"context"
	"errors"
	"micros/logging"
	"sync"

	"go.mongodb.org/mongo-driver/event"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type monitor struct {
	startedCommands map[int64]context.Context
	sync.RWMutex
}

const instrumentationName = "micros/database/mongo"

func (m *monitor) Started(ctx context.Context, e *event.CommandStartedEvent) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(semconv.PeerServiceKey.String("mongodb"))}

	ctx, _ = otel.Tracer(instrumentationName).Start(ctx, e.CommandName, opts...)

	logging.MongoStartedLog(ctx, e)

	m.Lock()
	m.startedCommands[e.RequestID] = ctx
	m.Unlock()
}

func (m *monitor) Succeeded(ctx context.Context, e *event.CommandSucceededEvent) {
	m.finished(e.CommandFinishedEvent, nil)
}

func (m *monitor) Failed(ctx context.Context, e *event.CommandFailedEvent) {
	m.finished(e.CommandFinishedEvent, errors.New(e.Failure))
}

func (m *monitor) finished(e event.CommandFinishedEvent, err error) {
	m.RLock()
	ctx, ok := m.startedCommands[e.RequestID]
	m.RUnlock()

	if ok {
		span := trace.SpanFromContext(ctx)
		defer span.End()

		m.Lock()
		delete(m.startedCommands, e.RequestID)
		m.Unlock()

		logging.MongoFinishedLog(ctx, e, err)
	}
}

func newMonitor() *event.CommandMonitor {
	m := monitor{
		startedCommands: make(map[int64]context.Context),
	}

	return &event.CommandMonitor{
		Started:   m.Started,
		Succeeded: m.Succeeded,
		Failed:    m.Failed,
	}
}
