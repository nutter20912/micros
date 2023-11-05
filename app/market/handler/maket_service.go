package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"micros/app/market/binance"
	"micros/event"
	marketV1 "micros/proto/market/v1"

	"go-micro.dev/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewMarketService(s micro.Service, e *event.Event) *MarketService {
	return &MarketService{Service: s, Event: e}
}

type MarketService struct {
	Service micro.Service

	Event *event.Event
}

func (s *MarketService) GetTradeStream(
	ctx context.Context,
	req *marketV1.GetTradeStreamResquest,
	stream marketV1.MarketService_GetTradeStreamStream,
) error {
	defer stream.Context().Done()

	if err := req.Validate(); err != nil {
		return stream.SendMsg(status.Error(codes.Internal, err.Error()))
	}

	// TODO Validate symbol
	c, err := binance.NewClient(context.Background()).Stream(*req.Symbol)
	if err != nil {
		return stream.SendMsg(status.Error(codes.Internal, err.Error()))
	}

	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Println("read message")
			rsp, err := c.ReadMessage()
			if err != nil {
				return stream.SendMsg(status.Error(codes.Internal, err.Error()))
			}

			var aggTradeMsg binance.AggTradeMessage
			json.Unmarshal([]byte(rsp), &aggTradeMsg)

			data := marketV1.AggTradeData{
				EventType:       aggTradeMsg.Data.EventType,
				EventTime:       aggTradeMsg.Data.EventTime,
				Price:           aggTradeMsg.Data.Price,
				Symbol:          aggTradeMsg.Data.Symbol,
				Quantity:        aggTradeMsg.Data.Quantity,
				TransactionTime: aggTradeMsg.Data.TransactionTime,
				IsSell:          aggTradeMsg.Data.IsSell}

			if data.Price == 0 {
				continue
			}

			if err := stream.Send(&marketV1.GetTradeStreamResponse{AggTrade: &data}); err != nil {
				return stream.SendMsg(status.Error(codes.Internal, err.Error()))
			}
		}
	}
}
