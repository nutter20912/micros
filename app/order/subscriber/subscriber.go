package subscriber

import (
	"micros/event"

	"go-micro.dev/v4"
)

func Register(s micro.Service) {
	walletSub := &walletSubscriber{Service: s}

	r := map[string]interface{}{
		event.WALLET_TRANSACTION: walletSub.addDepositOrderEvent,
	}

	for k, v := range r {
		micro.RegisterSubscriber(k, s.Server(), v)
	}
}
