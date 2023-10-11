package subscriber

import (
	"context"
	"micros/app/order/models"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"
)

type InsertOrderEvent struct{}

func (e *InsertOrderEvent) Handle(ctx context.Context, event *walletV1.TransactionEvent) error {
	depositOrder, err := new(models.DepositOrder).Get(event.OrderId)
	if err != nil {
		return err
	}

	status := orderV1.DepositStatus_DEPOSIT_STATUS_FAILED
	if event.Success {
		status = orderV1.DepositStatus_DEPOSIT_STATUS_COMPLETED
	}

	new(models.DepositOrderEvent).Add(depositOrder.Id, depositOrder.UserId, depositOrder.Amount, status)

	return nil
}
