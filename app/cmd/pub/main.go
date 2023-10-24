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
	defer fmt.Println("leave jetStream")
	nc, err := nats.Connect("nats")
	if err != nil {
		return nil, err
	}

	go func(ctx context.Context) {
		defer func() {
			nc.Close()
			fmt.Println("leave jetStreamFunc")
		}()

		<-ctx.Done()
	}(ctx)

	js, _ := jetstream.New(nc)

	return js, nil
}

func stream(ctx context.Context, js jetstream.JetStream) (jetstream.Stream, error) {
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      "market",
		Subjects:  []string{"market.*"},
		Retention: jetstream.InterestPolicy,
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
	js, err := jetStream(ctx)
	if err != nil {
		log.Fatalf("Connect Error: %v", err)
	}

	//_, err = stream(ctx, js)
	s, err := js.Stream(context.Background(), "market")
	if err != nil {
		log.Fatalf("Stream Error: %v", err)
	}

	go pub(ctx, js)
	go sub(ctx, s)

	time.Sleep(time.Second * 10)
}

func pub(ctx context.Context, js jetstream.JetStream) {
	js.Publish(ctx, "market.aaa", nil)
}

func sub(ctx context.Context, s jetstream.Stream) {
	defer fmt.Println("leave sub")
	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		FilterSubject: "market.test",
	})

	if err != nil {
		log.Fatalln(err)
	}

	cc, _ := c.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		fmt.Println("received msg on", msg.Subject())
	})

	defer cc.Stop()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
