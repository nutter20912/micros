package subscriber

import (
	"context"
	"fmt"
	"micros/app/market/event"
	"micros/database/redis"
	marketV1 "micros/proto/market/v1"
	orderV1 "micros/proto/order/v1"
	"strconv"

	"go-micro.dev/v4"
)

type orderSubscriber struct {
	Service micro.Service
}

// 搓合
func (s *orderSubscriber) matchOrder(ctx context.Context, e *orderV1.SpotCreatedEventMessage) error {
	price := e.Data.Price

	if e.Data.Type == orderV1.OrderType_ORDER_TYPE_MARKET {
		rdb := redis.Get()
		val, err := rdb.HGet(ctx, "price", e.Data.Symbol).Result()
		if err != nil {
			fmt.Println(err)
		}

		price, _ = strconv.ParseFloat(val, 64)
	}

	msg := &marketV1.OrderMatchedEventMessage{
		UserId:   e.Data.UserId,
		OrderId:  e.Data.OrderId,
		Symbol:   e.Data.Symbol,
		Price:    price,
		Quantity: e.Data.Quantity,
	}

	event.OrderMatched{Client: s.Service.Client()}.Dispatch(msg)
	return nil
}
