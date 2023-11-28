package subscriber

import (
	"micros/event"
	"micros/queue"

	"go-micro.dev/v4"
)

func Register(s micro.Service) {
	orderSub := &orderSubscriber{Service: s}

	r := map[string]interface{}{
		event.ORDER_SPOT_CREATED: orderSub.matchOrder,
	}

	queue.RegisterSubscriber(s.Server(), r)
}
