package natsjs

import (
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"
)

type optionsKey struct{}
type publishOptionsKey struct{}
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

func ConsumerConfig(cfg jetstream.ConsumerConfig) broker.SubscribeOption {
	return setSubscribeOption(consumerConfigKey{}, cfg)
}

func PublishOptions(opts ...jetstream.PublishOpt) client.PublishOption {
	return setPublishOptions(publishOptionsKey{}, opts)
}

// server.Subscriber
func SubscriberStreamConfig(cfg jetstream.StreamConfig) server.SubscriberOption {
	return setSubscriberOption(streamConfigKey{}, cfg)
}

func SubscriberConsumerConfig(cfg jetstream.ConsumerConfig) server.SubscriberOption {
	return setSubscriberOption(consumerConfigKey{}, cfg)
}
