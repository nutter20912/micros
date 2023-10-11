package event

import (
	"context"
	"fmt"
	orderV1 "micros/proto/order/v1"

	"github.com/spf13/viper"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type OrderCreated struct {
	Client client.Client
}

func (o OrderCreated) Topic() string {
	return viper.GetString("topic.order.created")
}

func (o OrderCreated) Dispatch(depositOrderEvent *orderV1.DepositOrderEvent) {
	pub := micro.NewEvent(o.Topic(), o.Client)

	msg := depositOrderEvent

	if err := pub.Publish(context.Background(), msg); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}
