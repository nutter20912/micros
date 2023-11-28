package queue

import (
	"micros/broker/natsjs"
	"strings"

	"github.com/nats-io/nats.go/jetstream"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

func RegisterSubscriber(server server.Server, r map[string]interface{}) {
	for k, v := range r {
		opt := natsjs.SubscriberConsumerConfig(jetstream.ConsumerConfig{
			Durable: strings.Replace(k, ".", "_", -1),
		})
		micro.RegisterSubscriber(k, server, v, opt)
	}
}
