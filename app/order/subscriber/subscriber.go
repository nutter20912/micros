package subscriber

import (
	"context"
	"micros/app/order/models"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"
)

type AddOrderEvent struct{}

func (e *AddOrderEvent) Handle(ctx context.Context, msg *walletV1.TransactionEventMessage) error {
	depositOrder, err := new(models.DepositOrder).Get(msg.OrderId)
	if err != nil {
		return err
	}

	status := orderV1.DepositStatus_DEPOSIT_STATUS_FAILED
	if msg.Success {
		status = orderV1.DepositStatus_DEPOSIT_STATUS_COMPLETED
	}

	new(models.DepositOrderEvent).Add(depositOrder.Id, depositOrder.UserId, depositOrder.Amount, status)

	return nil
}
