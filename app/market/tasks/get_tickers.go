package tasks

import (
	"context"
	"fmt"
	"log"
	"micros/app/market/binance"
	"micros/database/redis"

	"go-micro.dev/v4"
)

type getTickers struct {
	Service micro.Service
}

func (e *getTickers) Run() {
	rdb := redis.Get()
	stream, _ := binance.NewClient(context.Background()).MiniTickersStream()

	defer stream.Close()

	for {
		rsp, err := stream.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		tickers, ok := rsp.(*binance.MiniTickerArrMessage)
		if !ok {
			continue
		}

		if len(tickers.Data) == 0 {
			continue
		}

		price := map[string]interface{}{}
		for _, v := range tickers.Data {
			price[v.Symbol] = v.Close
		}

		if err := rdb.HSet(context.Background(), "price", price).Err(); err != nil {
			log.Println(err)
		}
	}
}
