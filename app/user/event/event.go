package event

import (
	"context"
	"fmt"
	"micros/event"
	userV1 "micros/proto/user/v1"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type UserCreated struct {
	Client client.Client
}

func (u UserCreated) Dispatch(userId string) {
	pub := micro.NewEvent(event.USER_CREATED, u.Client)

	msg := &userV1.RegisteredEventMessage{UserId: userId}

	if err := pub.Publish(context.Background(), msg); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}
