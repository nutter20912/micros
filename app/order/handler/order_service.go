package handler

import (
	"context"
	"fmt"

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

	depositOrderEvent, err := new(models.DepositOrderEvent).Create(userId, req.Amount)
	if err != nil {
		fmt.Println(err)
		return microErrors.BadRequest("123", "create fail")
	}

	rsp.Data = &orderV1.DepositOrderEvent{
		Id:     depositOrderEvent.Id.String(),
		UserId: depositOrderEvent.UserId,
		Status: depositOrderEvent.Status,
		Amount: depositOrderEvent.Amount,
	}

	event.OrderCreated{Client: s.Service.Client()}.Dispatch(rsp.Data)

	return nil
}

func (s *OrderService) GetDepositEvent(
	ctx context.Context,
	req *orderV1.GetDepositEventRequest,
	rsp *orderV1.GetDepositEventResponse,
) error {
	return nil
}

func (s *OrderService) GetDeposit(
	ctx context.Context,
	req *orderV1.GetDepositRequest,
	rsp *orderV1.GetDepositResponse,
) error {
	return nil
}
