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

func (o OrderCreated) Topic() string {
	return "order.deposit.created"
}

func (o OrderCreated) Dispatch(depositOrderEvent *orderV1.DepositOrderEvent, opts ...event.DispatchOption) error {
	pub := micro.NewEvent(o.Topic(), o.Client)

	msg := depositOrderEvent

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

type OrderCheck struct {
	Client client.Client
}

func (o OrderCheck) Topic() string {
	return "order.check"
}

func (o OrderCheck) Dispatch(m *orderV1.OrderCheckEventMessage, opts ...event.DispatchOption) error {
	pub := micro.NewEvent(o.Topic(), o.Client)

	mdOpts := map[string]string{}
	for _, o := range opts {
		o(mdOpts)
	}

	ctx := metadata.NewContext(context.Background(), mdOpts)

	if err := pub.Publish(ctx, m); err != nil {
		return fmt.Errorf("publish error: %v", err)
	}

	return nil
}
