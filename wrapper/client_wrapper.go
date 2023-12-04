package wrapper

import (
	"context"
	"micros/logging"

	"go-micro.dev/v4/client"
)

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	logging.PublishLog(ctx, l, p)

	return l.Client.Publish(ctx, p, opts...)
}

func NewClientWrapper(c client.Client) client.Client {
	return &logWrapper{c}
}
