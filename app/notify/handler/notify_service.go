package handler

import (
	"context"
	"fmt"

	"micros/broker/natsjs"
	"micros/event"
	notifyV1 "micros/proto/notify/v1"

	"github.com/nats-io/nats.go/jetstream"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type notifyService struct {
	Service micro.Service

	Event *event.Event
}

func NewNotifyService(s micro.Service, e *event.Event) *notifyService {
	return &notifyService{Service: s, Event: e}
}

func (ns *notifyService) GetStream(
	ctx context.Context,
	req *emptypb.Empty,
	stream notifyV1.NotifyService_GetStreamStream,
) error {
	defer stream.Context().Done()

	userId, _ := metadata.Get(ctx, "user_id")
	topic := fmt.Sprintf("notify.user.%v", userId)

	opts := natsjs.ConsumerConfig(jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})

	s, err := ns.Service.Options().Broker.Subscribe(topic, func(p broker.Event) error {
		msg := &notifyV1.GetStreamResponse{Data: &notifyV1.ChannelMessage{Payload: p.Message().Body}}

		if err := stream.Send(msg); err != nil {
			return stream.SendMsg(status.Error(codes.Internal, err.Error()))
		}

		return nil
	}, opts)

	if err != nil {
		return stream.SendMsg(status.Error(codes.Internal, err.Error()))
	}

	defer s.Unsubscribe()

	<-ctx.Done()

	return nil
}
