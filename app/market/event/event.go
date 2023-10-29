package event

import (
	"context"
	"fmt"
	"micros/event"
	marketV1 "micros/proto/market/v1"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type PriceChanged struct {
	Client client.Client
}

func (e PriceChanged) Dispatch(msg *walletV1.TransactionEventMessage) {
	pub := micro.NewEvent(event.MARKET_PRICE_CHANGED, e.Client)

	if err := pub.Publish(context.Background(), msg); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}

type OrderMatched struct {
	Client client.Client
}

func (e OrderMatched) Dispatch(msg *marketV1.OrderMatchedEventMessage) {
	pub := micro.NewEvent(event.MARKET_MATCHED, e.Client)

	if err := pub.Publish(context.Background(), msg); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}
