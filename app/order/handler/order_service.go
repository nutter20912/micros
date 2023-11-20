package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	OrderEvent "micros/app/order/event"
	"micros/app/order/models"
	"micros/database/mongo"
	"micros/event"
	orderV1 "micros/proto/order/v1"

	"go-micro.dev/v4"
	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	Service micro.Service

	Event *event.Event
}

func NewOrderService(s micro.Service, e *event.Event) *OrderService {
	return &OrderService{Service: s, Event: e}
}

func (s *OrderService) CreateDepositEvent(
	ctx context.Context,
	req *orderV1.CreateDepositEventRequest,
	rsp *orderV1.CreateDepositEventResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")

	depositOrderEvent, err := new(models.DepositOrderEvent).Create(userId, req.GetAmount())
	if err != nil {
		return microErrors.BadRequest("123", "create fail")
	}

	rsp.Data = &orderV1.DepositOrderEvent{
		Id:      depositOrderEvent.Id.Hex(),
		OrderId: depositOrderEvent.OrderId,
		UserId:  depositOrderEvent.UserId,
		Status:  depositOrderEvent.Status,
		Amount:  depositOrderEvent.Amount,
	}

	if err := s.Event.Dispatch(OrderEvent.DepositOrderCreated{Payload: rsp.Data}); err != nil {
		return microErrors.InternalServerError("123", "Dispatch error: %v", err)
	}

	return nil
}

func (s *OrderService) GetDepositEvent(
	ctx context.Context,
	req *orderV1.GetDepositEventRequest,
	rsp *orderV1.GetDepositEventResponse,
) error {
	events, err := new(models.DepositOrderEvent).Get(req.GetOrderId())
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data []*orderV1.DepositOrderEvent
	bytes, _ := json.Marshal(events)
	json.Unmarshal(bytes, &data)

	rsp.Data = data

	return nil
}

func (s *OrderService) GetDeposit(
	ctx context.Context,
	req *orderV1.GetDepositRequest,
	rsp *orderV1.GetDepositResponse,
) error {
	order, err := new(models.DepositOrder).Get(req.GetId())
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data *orderV1.DepositOrder
	bytes, _ := json.Marshal(order)
	json.Unmarshal(bytes, &data)

	rsp.Data = data

	return nil
}

func (s *OrderService) CreateSpotEvent(
	ctx context.Context,
	req *orderV1.CreateSpotEventRequest,
	rsp *orderV1.CreateSpotEventResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")

	spotOrderEvent := models.SpotOrderEvent{
		UserId:   userId,
		Symbol:   req.Symbol,
		Quantity: req.Quantity,
		Side:     req.Side,
		Type:     req.Type,
		Status:   orderV1.SpotStatus_SPOT_STATUS_NEW,
	}

	if req.Type == orderV1.OrderType_ORDER_TYPE_LIMIT {
		if req.Price == nil {
			return microErrors.BadRequest("222", "LIMIT must required Price")
		}

		spotOrderEvent.Price = *req.Price
	}

	if err := spotOrderEvent.Create(); err != nil {
		return microErrors.BadRequest("123", "create fail")
	}

	var data *orderV1.SpotOrderEvent
	bytes, _ := json.Marshal(spotOrderEvent)
	json.Unmarshal(bytes, &data)

	rsp.Data = data

	if err := s.Event.Dispatch(OrderEvent.SpotOrderCreated{Payload: rsp.Data}); err != nil {
		return microErrors.InternalServerError("123", "Dispatch error: %v", err)
	}

	s.Event.Dispatch(event.Notify{
		Channel: fmt.Sprintf("user.%s", spotOrderEvent.UserId),
		Name:    "SpotOrderEvent",
		Payload: spotOrderEvent})

	return nil
}

func (s *OrderService) GetSpotEvent(
	ctx context.Context,
	req *orderV1.GetSpotEventRequest,
	rsp *orderV1.GetSpotEventResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	startDate, err := time.Parse(time.RFC3339, req.GetStartDate())
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}
	endDate, err := time.Parse(time.RFC3339, req.GetEndDate())
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")

	events, paginator, err := new(models.SpotOrderEvent).Get(
		req.Page,
		req.Limit,
		mongo.FilterField("user_id", userId),
		mongo.FilterDateRange("time", startDate, endDate))
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data []*orderV1.SpotOrderEvent
	bytes, _ := json.Marshal(events)
	json.Unmarshal(bytes, &data)

	var p *orderV1.Paginator
	bytes, _ = json.Marshal(paginator)
	json.Unmarshal(bytes, &p)

	rsp.Data = data
	rsp.Paginator = p

	return nil
}

func (s *OrderService) GetSpotPosition(
	ctx context.Context,
	req *orderV1.GetSpotPositionRequest,
	rsp *orderV1.GetSpotPositionResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")
	symbol := strings.ToUpper(req.GetSymbol())

	spotPositions, paginator, err := new(models.SpotPosition).Get(userId, symbol, req.Page, req.Limit)
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data []*orderV1.SpotPosition
	bytes, _ := json.Marshal(spotPositions)
	json.Unmarshal(bytes, &data)

	var p *orderV1.Paginator
	bytes, _ = json.Marshal(paginator)
	json.Unmarshal(bytes, &p)

	rsp.Data = data
	rsp.Paginator = p

	return nil
}

func (s *OrderService) GetSpotPositionClosed(
	ctx context.Context,
	req *orderV1.GetSpotPositionClosedRequest,
	rsp *orderV1.GetSpotPositionClosedResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")
	symbol := strings.ToUpper(req.GetSymbol())

	spotPositions, paginator, err := new(models.SpotPositionClosed).Get(userId, symbol, req.Page, req.Limit)
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data []*orderV1.SpotPositionClosed
	bytes, _ := json.Marshal(spotPositions)
	json.Unmarshal(bytes, &data)

	var p *orderV1.Paginator
	bytes, _ = json.Marshal(paginator)
	json.Unmarshal(bytes, &p)

	rsp.Data = data
	rsp.Paginator = p

	return nil
}

func (s *OrderService) GetPositionStream(
	ctx context.Context,
	req *orderV1.GetPositionStreamResquest,
	stream orderV1.OrderService_GetPositionStreamStream,
) error {
	defer stream.Context().Done()

	if err := req.Validate(); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")
	symbol := strings.ToUpper(req.GetSymbol())

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Second)

			var sp models.SpotPosition
			spotPositions, err := sp.GetList(userId, symbol)
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			var open []*orderV1.SpotPosition
			openBytes, _ := json.Marshal(spotPositions)
			json.Unmarshal(openBytes, &open)

			if err := stream.Send(&orderV1.GetPositionStreamResponse{Open: open}); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
}
