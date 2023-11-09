package subscriber

import (
	"context"
	"errors"
	"fmt"
	"micros/app/order/models"
	"micros/event"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"
	"micros/queue"

	"go-micro.dev/v4"
	"go-micro.dev/v4/codec/proto"
)

type walletSubscriber struct {
	Service micro.Service

	Event *event.Event
}

func (s *walletSubscriber) addOrderEvent(ctx context.Context, m *[]byte) error {
	microId, err := queue.MicroId(ctx)
	if err != nil {
		return err
	}

	msg := walletV1.TransactionEventMessage{}
	if err := new(proto.Marshaler).Unmarshal(*m, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case walletV1.WalletEventType_WALLET_EVENT_TYPE_DEPOSIT:
		return s.addDepositOrderEvent(ctx, microId, &msg)
	default:
		return fmt.Errorf("wrong trans type")
	}
}

func (s *walletSubscriber) addDepositOrderEvent(
	ctx context.Context,
	microId string,
	msg *walletV1.TransactionEventMessage,
) error {
	if err := queue.CheckMsgId(new(models.DepositOrderEvent), microId); err != nil {
		return err
	}

	depositOrder, err := new(models.DepositOrder).Get(msg.OrderId)
	if err != nil {
		return errors.New("deposit_order not found")
	}

	status := orderV1.DepositStatus_DEPOSIT_STATUS_FAILED
	if msg.Success {
		status = orderV1.DepositStatus_DEPOSIT_STATUS_COMPLETED
	}

	event := &models.DepositOrderEvent{
		MsgId:   microId,
		OrderId: depositOrder.Id,
		UserId:  depositOrder.UserId,
		Amount:  depositOrder.Amount,
		Status:  status,
		Memo:    msg.Memo,
	}

	new(models.DepositOrderEvent).Add(event)

	return nil
}

func (s *walletSubscriber) addSpotOrderEvent(
	ctx context.Context,
	msg *walletV1.BalanceCheckedEventMessage,
) error {
	microId, err := queue.MicroId(ctx)
	if err != nil {
		return err
	}

	if err := queue.CheckMsgId(new(models.SpotOrderEvent), microId); err != nil {
		return err
	}

	spotOrderEvent := models.SpotOrderEvent{OrderId: msg.OrderId}

	count, err := spotOrderEvent.Count()
	if err != nil {
		return err
	}

	if count > 1 {
		return queue.ErrMessageConflicted
	}

	if err := spotOrderEvent.Last(); err != nil {
		return err
	}

	spotOrderEvent.MsgId = microId

	if msg.Success {
		spotOrderEvent.Price = msg.Price
		spotOrderEvent.Status = orderV1.SpotStatus_SPOT_STATUS_FILLED
	} else {
		spotOrderEvent.Status = orderV1.SpotStatus_SPOT_STATUS_REJECTED
		spotOrderEvent.Memo = msg.Memo
	}

	if err := spotOrderEvent.Add(); err != nil {
		return err
	}

	s.Event.Dispatch(event.Notify{
		Channel: fmt.Sprintf("user.%s", spotOrderEvent.UserId),
		Name:    "SpotOrderEvent",
		Payload: spotOrderEvent})

	p := models.SpotPosition{
		UserId:   spotOrderEvent.UserId,
		OrderId:  spotOrderEvent.OrderId,
		Symbol:   spotOrderEvent.Symbol,
		Side:     spotOrderEvent.Side,
		Price:    spotOrderEvent.Price,
		Quantity: spotOrderEvent.Quantity,
	}

	p.Upsert()

	return nil
}
