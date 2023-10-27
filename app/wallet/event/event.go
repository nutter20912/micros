package event

import (
	"context"
	"fmt"
	walletV1 "micros/proto/wallet/v1"

	"github.com/spf13/viper"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
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
