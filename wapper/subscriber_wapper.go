package wapper

import (
	"context"
	"log"
	"micros/event"

	"go-micro.dev/v4/server"
)

func LogSubWapper() server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			if err := event.Validate(ctx); err != nil {
				return event.ErrReportOrIgnore(err)
			}

			log.Printf("[sub_log] topic: %v", msg.Topic())
			log.Printf("[sub_log] payload: %v", msg.Payload())

			if err := fn(ctx, msg); err != nil {
				return event.ErrReportOrIgnore(err)
			}

			return nil
		}
	}
}
