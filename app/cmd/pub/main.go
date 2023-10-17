package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func jetStream(ctx context.Context) (jetstream.JetStream, error) {
	nc, err := nats.Connect("nats")
	if err != nil {
		return nil, err
	}

	defer nc.Close()

	js, _ := jetstream.New(nc)

	return js, nil
}

func stream(ctx context.Context, js jetstream.JetStream) (jetstream.Stream, error) {
	//js.DeleteStream(ctx, "EVENTS")
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      "EVENTS",
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"events.>"},
	})

	if err != nil {
		return nil, err
	}

	i, _ := s.Info(ctx)
	ib, _ := json.Marshal(i)
	fmt.Println(string(ib))

	return s, nil
}

func main() {
	ctx := context.Background()
	js, _ := jetStream(ctx)
	s, _ := stream(ctx, js)

	go pub(ctx, js)
	go sub(ctx, s)

	time.Sleep(time.Minute)
}

func pub(ctx context.Context, js jetstream.JetStream) {
	js.Publish(ctx, "events.us.page_loaded", nil)
}

func sub(ctx context.Context, s jetstream.Stream) {
	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
	})

	if err != nil {
		log.Fatalln(err)
	}

	cc, _ := c.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		fmt.Println("received msg on", msg.Subject())
	})

	defer cc.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
