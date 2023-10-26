package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	OrderAction "micros/app/order/action"
	OrderEvent "micros/app/order/event"
	"micros/app/order/models"
	"micros/event"
	orderV1 "micros/proto/order/v1"

	"github.com/oklog/ulid/v2"
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

	err = OrderEvent.OrderCreated{Client: s.Service.Client()}.Dispatch(rsp.Data)
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

func (s *OrderService) CreateSpot(
	ctx context.Context,
	req *orderV1.CreateSpotRequest,
	rsp *orderV1.CreateSpotResponse,
) error {
	userId, _ := metadata.Get(ctx, "user_id")
	callbackSubject := fmt.Sprintf("order.callback.%s", ulid.Make().String())

	ttl := time.Second * 5

	result := []*orderV1.CheckCallbackMessage{}
	var wg sync.WaitGroup
	wg.Add(1)

	go OrderAction.WaitForCheck(
		s.Service.Options().Broker, &result, &wg, callbackSubject, ttl)

	msg := orderV1.OrderCheckEventMessage{
		UserId:          userId,
		CallbackSubject: callbackSubject}
	oe := OrderEvent.OrderCheck{Client: s.Service.Client()}

	if err := oe.Dispatch(&msg, event.SetTTL(ttl)); err != nil {
		return microErrors.InternalServerError("222", err.Error())
	}

	wg.Wait()

	for _, v := range result {
		if !v.Success {
			return microErrors.BadRequest("222", v.Msg)
		}
	}

	return nil
}
