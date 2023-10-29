package subscriber

import (
	"micros/event"

	"go-micro.dev/v4"
)

func Register(s micro.Service) {
	orderSub := &orderSubscriber{Service: s}

	r := map[string]interface{}{
		event.ORDER_SPOT_CREATED: orderSub.matchOrder,
	}

	for k, v := range r {
		micro.RegisterSubscriber(k, s.Server(), v)
	}
}
