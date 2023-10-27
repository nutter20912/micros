package action

import (
	"context"
	"fmt"
	"micros/broker/natsjs"
	orderV1 "micros/proto/order/v1"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/codec/proto"
)

func WaitForCheck(
	b broker.Broker,
	res *[]*orderV1.CheckCallbackMessage,
	wg *sync.WaitGroup,
	callbackSubject string,
	ttl time.Duration,
) {
	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer func() {
		cancel()
		wg.Done()
	}()

	opts := natsjs.ConsumerConfig(jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverLastPolicy,
	})

	s, err := b.Subscribe(callbackSubject, func(p broker.Event) error {
		var wap []byte
		b.Options().Codec.Unmarshal(p.Message().Body, &wap)

		var msg orderV1.CheckCallbackMessage
		marshaler := proto.Marshaler{}
		if err := marshaler.Unmarshal(wap, &msg); err != nil {
			return err
		}

		*res = append(*res, &msg)

		return nil
	}, opts)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer s.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if len(*res) >= 2 {
			return
		}

		time.Sleep(time.Millisecond * 100)
	}

	//switch v.Data.(type) {
	//case *orderV1.CheckCallbackMessage_Market:
	//case *orderV1.CheckCallbackMessage_Wallet:
	//}

}
