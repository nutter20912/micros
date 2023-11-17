package event

import (
	"context"
	"micros/event"
	userV1 "micros/proto/user/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type UserCreated struct {
	Payload string
}

func (o UserCreated) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(event.USER_CREATED, c)

	msg := &userV1.RegisteredEventMessage{UserId: o.Payload}

	return pub.Publish(ctx, msg)
}
