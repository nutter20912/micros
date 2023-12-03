package wrapper

import (
	"context"
	"micros/logging"
	"micros/queue"

	"go-micro.dev/v4/server"
)

func LogSubWrapper() server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			if err := queue.Validate(ctx); err != nil {
				return queue.ErrReportOrIgnore(err)
			}

			logging.EventLog(ctx, msg)

			if err := fn(ctx, msg); err != nil {
				return queue.ErrReportOrIgnore(err)
			}

			return nil
		}
	}
}
