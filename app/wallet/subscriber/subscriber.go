package subscriber

import (
	"micros/event"
	"micros/queue"

	"go-micro.dev/v4"
)

func Register(s micro.Service, e *event.Event) {
	orderSub := &orderSubscriber{Service: s, Event: e}
	userSub := &userSubscriber{Service: s}

	r := map[string]interface{}{
		event.ORDER_DEPOSIT_CREATED: orderSub.addEventByDeposit,
		event.MARKET_MATCHED:        orderSub.checkBalance,

		event.USER_CREATED: userSub.initWalletEvent,
	}

	queue.RegisterSubscriber(s.Server(), r)
}
