package subscriber

import (
	"context"
	"errors"
	"micros/app/order/models"
	baseEvent "micros/event"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
)

type walletSubscriber struct {
	Service micro.Service
}

func (s *walletSubscriber) addDepositOrderEvent(ctx context.Context, msg *walletV1.TransactionEventMessage) error {
	microId, err := baseEvent.MicroId(ctx)
	if err != nil {
		return err
	}

	if err := validate(ctx, microId); err != nil {
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
