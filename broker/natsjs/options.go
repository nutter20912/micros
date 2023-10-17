package natsjs

import (
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/server"
)

type optionsKey struct{}
type drainConnectionKey struct{}
type streamConfigKey struct{}
type consumerConfigKey struct{}

// Options accepts nats.Options.
func Options(opts nats.Options) broker.Option {
	return setBrokerOption(optionsKey{}, opts)
}

// DrainConnection will drain subscription on close.
func DrainConnection() broker.Option {
	return setBrokerOption(drainConnectionKey{}, struct{}{})
}

func StreamConfig(cfg jetstream.StreamConfig) server.SubscriberOption {
	return setSubscriberOption(streamConfigKey{}, cfg)
}

func ConsumerConfig(cfg jetstream.ConsumerConfig) server.SubscriberOption {
	return setSubscriberOption(consumerConfigKey{}, cfg)
}
