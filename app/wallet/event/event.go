package event

import (
	"context"
	"fmt"
	"micros/event"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/codec/proto"
)

type TransactionEvent struct {
	Payload *walletV1.TransactionEventMessage
}

func (e TransactionEvent) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(event.WALLET_TRANSACTION, c)

	marshaler := proto.Marshaler{}
	b, err := marshaler.Marshal(e.Payload)
	if err != nil {
		return fmt.Errorf("error Marshal: %v", err)
	}

	return pub.Publish(ctx, b)
}

type BalanceChecked struct {
	Payload *walletV1.BalanceCheckedEventMessage
}

func (e BalanceChecked) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(event.WALLET_BALANCE_CHECKED, c)

	return pub.Publish(ctx, e.Payload)
}
