package subscriber

import (
	"micros/event"

	"go-micro.dev/v4"
)

func Register(s micro.Service) {
	orderSub := &orderSubscriber{Service: s}
	userSub := &userSubscriber{Service: s}

	r := map[string]interface{}{
		event.ORDER_DEPOSIT_CREATED: orderSub.addWalletEvent,

		event.USER_CREATED: userSub.initWalletEvent,
	}

	for k, v := range r {
		micro.RegisterSubscriber(k, s.Server(), v)
	}
}
