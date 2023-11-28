package subscriber

import (
	"micros/event"
	"micros/queue"

	"go-micro.dev/v4"
)

func Register(s micro.Service, e *event.Event) {
	walletSub := &walletSubscriber{Service: s, Event: e}

	r := map[string]interface{}{
		event.WALLET_TRANSACTION:     walletSub.addOrderEvent,
		event.WALLET_BALANCE_CHECKED: walletSub.addSpotOrderEvent,
	}

	queue.RegisterSubscriber(s.Server(), r)
}
