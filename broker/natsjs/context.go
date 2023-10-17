package natsjs

import (
	"context"

	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/server"
)

// setBrokerOption returns a function to setup a context with given value.
func setBrokerOption(k, v interface{}) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

func setSubscriberOption(k, v interface{}) server.SubscriberOption {
	return func(o *server.SubscriberOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
