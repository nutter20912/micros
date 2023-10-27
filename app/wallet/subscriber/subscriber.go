package subscriber

import "go-micro.dev/v4"

func Register(s micro.Service) {
	orderSub := &orderSubscriber{Service: s}
	userSub := &userSubscriber{Service: s}

	r := map[string]interface{}{
		"order.deposit.created": orderSub.depositCreated,
		"order.spot.created":    orderSub.spotCreated,

		"user.registered": userSub.registered,
	}

	for k, v := range r {
		micro.RegisterSubscriber(k, s.Server(), v)
	}
}
