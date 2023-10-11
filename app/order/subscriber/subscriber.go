package subscriber

import (
	"context"
	"micros/app/order/models"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"
)

type InsertOrderEvent struct{}

func (e *InsertOrderEvent) Handle(ctx context.Context, event *walletV1.TransactionEvent) error {
	depositOrder, err := models.GetDepositOrder(context.Background(), event.OrderId)
	if err != nil {
		return err
	}

	status := orderV1.DepositStatus_DEPOSIT_STATUS_FAILED
	if event.Success {
		status = orderV1.DepositStatus_DEPOSIT_STATUS_COMPLETED
	}

	models.InertDepositOrderEvent(depositOrder.Id, depositOrder.UserId, depositOrder.Amount, status)

	return nil
}
