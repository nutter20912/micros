package wapper

import (
	"context"
	"log"
	"micros/queue"

	"go-micro.dev/v4/server"
)

func LogSubWapper() server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			if err := queue.Validate(ctx); err != nil {
				return queue.ErrReportOrIgnore(err)
			}

			log.Printf("[sub_log] topic: %v", msg.Topic())
			log.Printf("[sub_log] payload: %v", msg.Payload())

			if err := fn(ctx, msg); err != nil {
				return queue.ErrReportOrIgnore(err)
			}

			return nil
		}
	}
}
