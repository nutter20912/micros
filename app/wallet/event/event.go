package event

import (
	"context"
	"fmt"
	"micros/event"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type TransactionEvent struct {
	Client client.Client
}

func (e TransactionEvent) Dispatch(msg *walletV1.TransactionEventMessage) {
	pub := micro.NewEvent(event.WALLET_TRANSACTION, e.Client)

	if err := pub.Publish(context.Background(), msg); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}
