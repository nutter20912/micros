package handler

import (
	"context"
	"encoding/json"
	"time"

	"micros/app/wallet/models"
	"micros/event"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewWalletService(s micro.Service, e *event.Event) *WalletService {
	return &WalletService{Service: s, Event: e}
}

type WalletService struct {
	Service micro.Service

	Event *event.Event
}

func (s *WalletService) GetEvent(
	ctx context.Context,
	req *walletV1.GetEventRequest,
	rsp *walletV1.GetEventResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")

	events, paginatior, err := new(models.WalletEvent).Get(userId, req.Page, req.Limit)
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data []*walletV1.WalletEvent
	bytes, _ := json.Marshal(events)
	json.Unmarshal(bytes, &data)

	var p *walletV1.Paginator
	bytes, _ = json.Marshal(paginatior)
	json.Unmarshal(bytes, &p)

	rsp.Data = data
	rsp.Paginator = p

	return nil
}

func (s *WalletService) Get(
	ctx context.Context,
	req *emptypb.Empty,
	rsp *walletV1.GetResponse,
) error {
	userId, _ := metadata.Get(ctx, "user_id")

	wallet, err := new(models.Wallet).Get(userId)
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data *walletV1.Wallet
	bytes, _ := json.Marshal(wallet)
	json.Unmarshal(bytes, &data)

	rsp.Data = data

	return nil
}

func (s *WalletService) GetWalletStream(
	ctx context.Context,
	req *walletV1.GetWalletStreamResquest,
	stream walletV1.WalletService_GetWalletStreamStream,
) error {
	defer stream.Context().Done()

	userId, _ := metadata.Get(ctx, "user_id")

	eventCursor := req.EventCursor

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Second)

			wallet, err := new(models.Wallet).Get(userId)
			if err != nil {
				return stream.SendMsg(status.Error(codes.Internal, err.Error()))
			}
			var info *walletV1.Wallet
			infoBytes, _ := json.Marshal(wallet)
			json.Unmarshal(infoBytes, &info)

			walletEvents, err := new(models.WalletEvent).GetEvents(userId, eventCursor)
			if err != nil {
				return stream.SendMsg(status.Error(codes.Internal, err.Error()))
			}

			if len(walletEvents) != 0 {
				lastId := walletEvents[0].Id.Hex()
				eventCursor = &lastId
			}

			var events []*walletV1.WalletEvent
			walletEventsBytes, _ := json.Marshal(walletEvents)
			json.Unmarshal(walletEventsBytes, &events)

			if err := stream.Send(&walletV1.GetWalletStreamResponse{Info: info, Events: events}); err != nil {
				return stream.SendMsg(status.Error(codes.Internal, err.Error()))
			}
		}
	}
}
