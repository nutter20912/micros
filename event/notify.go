package event

import (
	"context"
	"fmt"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

type Notify struct {
	Channel string
	Name    string
	Payload interface{}
}

func (n Notify) Publish(ctx context.Context, c client.Client) error {
	pub := micro.NewEvent(fmt.Sprintf("notify.%s", n.Channel), c)

	return pub.Publish(ctx, n)
}
