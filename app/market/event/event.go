package event

import (
	"context"
	"micros/event"
	marketV1 "micros/proto/market/v1"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type PriceChanged struct {
	Client  client.Client
	Payload *walletV1.TransactionEventMessage
}

func (e PriceChanged) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(event.MARKET_PRICE_CHANGED, c)

	return pub.Publish(ctx, e.Payload)
}

type OrderMatched struct {
	Payload *marketV1.OrderMatchedEventMessage
}

func (e OrderMatched) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(event.MARKET_MATCHED, c)

	return pub.Publish(ctx, e.Payload)
}
