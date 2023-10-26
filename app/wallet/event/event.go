package event

import (
	"context"
	"fmt"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"

	"github.com/spf13/viper"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/codec/proto"
)

type TransactionEvent struct {
	Client client.Client
}

func (e TransactionEvent) Topic() string {
	return viper.GetString("topic.wallet.transaction")
}

func (e TransactionEvent) Dispatch(msg *walletV1.TransactionEventMessage) {
	pub := micro.NewEvent(e.Topic(), e.Client)

	if err := pub.Publish(context.Background(), msg); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}

type CheckCallbackEvent struct {
	Client client.Client
}

func (e CheckCallbackEvent) Dispatch(topic string, msg *orderV1.CheckCallbackMessage) {
	pub := micro.NewEvent(topic, e.Client)

	p := proto.Marshaler{}
	b, _ := p.Marshal(msg)

	if err := pub.Publish(context.Background(), b); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}
