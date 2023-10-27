package subscriber

import "go-micro.dev/v4"

func Register(s micro.Service) {
	walletSub := &walletSubscriber{Service: s}

	r := map[string]interface{}{
		"wallet.transaction": walletSub.transactionEvent,
	}

	for k, v := range r {
		micro.RegisterSubscriber(k, s.Server(), v)
	}
}
