package handler

import (
	"context"
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
			rsp, err := c.ReadMessage()
			if err != nil {
				return stream.SendMsg(status.Error(codes.Internal, err.Error()))
			}

			streamRsp := &marketV1.GetTradeStreamResponse{}

			switch val := rsp.(type) {
			case *binance.AggTradeMessage:
				streamRsp.Data = &marketV1.GetTradeStreamResponse_AggTradeData{
					AggTradeData: &marketV1.AggTradeData{
						EventType:       val.Data.EventType,
						EventTime:       val.Data.EventTime,
						Price:           val.Data.Price,
						Symbol:          val.Data.Symbol,
						Quantity:        val.Data.Quantity,
						TransactionTime: val.Data.TransactionTime,
						IsSell:          val.Data.IsSell}}

			case *binance.KlineMessage:
				streamRsp.Data = &marketV1.GetTradeStreamResponse_KlineData{
					KlineData: &marketV1.KlineData{
						EventType: val.Data.EventType,
						EventTime: val.Data.EventTime,
						Symbol:    val.Data.Symbol,
						Kline: &marketV1.Kline{
							StartTime: val.Data.Kline.StartTime,
							EndTime:   val.Data.Kline.EndTime,
							Symbol:    val.Data.Kline.Symbol,
							Interval:  val.Data.Kline.Interval,
							Open:      val.Data.Kline.Open,
							Close:     val.Data.Kline.Close,
							High:      val.Data.Kline.High,
							Low:       val.Data.Kline.Low,
						},
					}}

			case *binance.DepthMessage:
				depthData := &marketV1.DepthData{
					Bids: []*marketV1.DepthData_DepthArr{},
					Asks: []*marketV1.DepthData_DepthArr{}}

				for _, v := range val.Data.Bids {
					depthData.Bids = append(depthData.Bids, &marketV1.DepthData_DepthArr{PriceAndQty: v})
				}

				for _, v := range val.Data.Asks {
					depthData.Asks = append(depthData.Asks, &marketV1.DepthData_DepthArr{PriceAndQty: v})
				}

				streamRsp.Data = &marketV1.GetTradeStreamResponse_DepthData{DepthData: depthData}

			default:
				continue
			}

			if err := stream.Send(streamRsp); err != nil {
				return stream.SendMsg(status.Error(codes.Internal, err.Error()))
			}
		}
	}
}
