package event

import (
	"context"
	"micros/event"
	orderV1 "micros/proto/order/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type DepositOrderCreated struct {
	Payload *orderV1.DepositOrderEvent
}

func (o DepositOrderCreated) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(event.ORDER_DEPOSIT_CREATED, c)

	msg := &orderV1.DepositCreatedEventMessage{
		UserId:  o.Payload.UserId,
		OrderId: o.Payload.OrderId,
		Amount:  o.Payload.Amount,
	}

	return pub.Publish(ctx, msg)
}

type SpotOrderCreated struct {
	Payload *orderV1.SpotOrderEvent
}

func (o SpotOrderCreated) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(event.ORDER_SPOT_CREATED, c)

	msg := &orderV1.SpotCreatedEventMessage{Data: o.Payload}

	return pub.Publish(ctx, msg)
}
