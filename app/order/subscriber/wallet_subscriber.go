package subscriber

import (
	"context"
	"errors"
	"fmt"
	"micros/app/order/models"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"
	"micros/queue"

	"go-micro.dev/v4"
	"go-micro.dev/v4/codec/proto"
)

type walletSubscriber struct {
	Service micro.Service
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
