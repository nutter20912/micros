package handler

import (
	"context"
	"encoding/json"

	"micros/app/order/event"
	"micros/app/order/models"
	orderV1 "micros/proto/order/v1"

	"go-micro.dev/v4"
	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
)

type OrderService struct {
	Service micro.Service
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

	err = event.OrderCreated{Client: s.Service.Client()}.Dispatch(rsp.Data)
	if err != nil {
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
