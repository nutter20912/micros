package event

import (
	"context"
	"fmt"
	"micros/event"
	orderV1 "micros/proto/order/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
)

type OrderCreated struct {
	Client client.Client
}

func (o OrderCreated) Dispatch(d *orderV1.DepositOrderEvent, opts ...event.DispatchOption) error {
	pub := micro.NewEvent(event.ORDER_DEPOSIT_CREATED, o.Client)

	msg := &orderV1.DepositCreatedEventMessage{
		UserId:  d.UserId,
		OrderId: d.OrderId,
		Amount:  d.Amount,
	}

	mdOpts := map[string]string{}
	for _, o := range opts {
		o(mdOpts)
	}

	ctx := metadata.NewContext(context.Background(), mdOpts)

	if err := pub.Publish(ctx, msg); err != nil {
		return fmt.Errorf("publish error: %v", err)
	}

	return nil
}
