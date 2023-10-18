package subscriber

import (
	"context"
	"errors"
	"micros/app/wallet/event"
	"micros/app/wallet/models"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
)

type OrderSubscriber struct {
	Service micro.Service
}

func (o *OrderSubscriber) DepositCreated(ctx context.Context, e *orderV1.DepositOrderEvent) error {
	msg := &walletV1.TransactionEventMessage{
		UserId:  e.UserId,
		OrderId: e.OrderId,
		Type:    walletV1.WalletEventType_WALLET_EVENT_TYPE_DEPOSIT,
		Success: true}

	err := func() error {
		_, err := new(models.Wallet).Get(e.UserId)
		if err != nil {
			return errors.New("wallet not found")
		}

		newWalletEvent := &walletV1.WalletEvent{
			UserId:  e.UserId,
			OrderId: e.OrderId,
			Change:  e.Amount,
			Type:    walletV1.WalletEventType_WALLET_EVENT_TYPE_DEPOSIT}

		if err := new(models.WalletEvent).Add(newWalletEvent); err != nil {
			return errors.New("add wallet_event error")
		}

		return nil
	}()

	if err != nil {
		msg.Memo = err.Error()
		msg.Success = false
	}

	event.TransactionEvent{Client: o.Service.Client()}.Dispatch(msg)

	return nil
}
