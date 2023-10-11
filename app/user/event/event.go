package event

import (
	"context"
	"fmt"
	userV1 "micros/proto/user/v1"

	"github.com/spf13/viper"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type UserCreated struct {
	Client client.Client
}

func (u UserCreated) Topic() string {
	return viper.GetString("topic.user.created")
}

func (u UserCreated) Dispatch(userId string) {
	pub := micro.NewEvent(u.Topic(), u.Client)

	msg := &userV1.RegisteredEvent{UserId: userId}

	if err := pub.Publish(context.Background(), msg); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}
