package subscriber

import (
	"context"
	"errors"
	"micros/app/order/models"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"
)

type WalletSubscriber struct{}

func (s *WalletSubscriber) TransactionEvent(ctx context.Context, msg *walletV1.TransactionEventMessage) error {
	depositOrder, err := new(models.DepositOrder).Get(msg.OrderId)
	if err != nil {
		return errors.New("deposit_order not found")
	}

	status := orderV1.DepositStatus_DEPOSIT_STATUS_FAILED
	if msg.Success {
		status = orderV1.DepositStatus_DEPOSIT_STATUS_COMPLETED
	}

	event := &models.DepositOrderEvent{
		OrderId: depositOrder.Id,
		UserId:  depositOrder.UserId,
		Amount:  depositOrder.Amount,
		Status:  status,
		Memo:    msg.Memo,
	}

	new(models.DepositOrderEvent).Add(event)

	return nil
}
