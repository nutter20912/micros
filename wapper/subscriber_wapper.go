package wapper

import (
	"context"
	"log"

	"go-micro.dev/v4/server"
)

func LogSubWapper() server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			log.Printf("[sub_log] topic: %v", msg.Topic())
			log.Printf("[sub_log] payload: %v", msg.Payload())
			err := fn(ctx, msg)

			return err
		}
	}
}
