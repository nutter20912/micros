package subscriber

import (
	"context"
	"fmt"
	walletEvent "micros/app/wallet/event"
	"micros/app/wallet/models"
	"micros/event"
	marketV1 "micros/proto/market/v1"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"
	"micros/queue"

	"go-micro.dev/v4"
)

type orderSubscriber struct {
	Service micro.Service

	Event *event.Event
}

func (o *orderSubscriber) addEventByDeposit(ctx context.Context, e *orderV1.DepositCreatedEventMessage) error {
	microId, err := queue.MicroId(ctx)
	if err != nil {
		return err
	}

	if err := queue.CheckMsgId(new(models.WalletEvent), microId); err != nil {
		return err
	}

	msg := &walletV1.TransactionEventMessage{}
	newEvent := &models.WalletEvent{
		MsgId:   microId,
		UserId:  e.UserId,
		OrderId: e.OrderId,
		Change:  e.Amount,
		Type:    walletV1.WalletEventType_WALLET_EVENT_TYPE_DEPOSIT}

	err = func() error {
		if _, err := new(models.Wallet).Get(newEvent.UserId); err != nil {
			return fmt.Errorf("wallet not found: %v", err)
		}

		if err := newEvent.Add(); err != nil {
			return fmt.Errorf("add wallet_event error: %v", err)
		}

		return nil
	}()

	if err != nil {
		msg.Success = false
		msg.Memo = err.Error()
	} else {
		msg.UserId = newEvent.UserId
		msg.OrderId = newEvent.OrderId
		msg.Type = walletV1.WalletEventType_WALLET_EVENT_TYPE_DEPOSIT
	}

	o.Event.Dispatch(walletEvent.TransactionEvent{Payload: msg})

	return nil
}

func (o *orderSubscriber) checkBalance(ctx context.Context, e *marketV1.OrderMatchedEventMessage) error {
	msg := &walletV1.BalanceCheckedEventMessage{
		UserId:   e.UserId,
		OrderId:  e.OrderId,
		Symbol:   e.Symbol,
		Price:    e.Price,
		Quantity: e.Quantity,
		Success:  true}

	wallet, err := new(models.Wallet).Get(e.UserId)
	if err != nil {
		return fmt.Errorf("wallet not found: %v", err)
	}

	if wallet.Amount < 10 {
		msg.Success = false
		msg.Memo = "balancce not enough"
	}

	o.Event.Dispatch(walletEvent.BalanceChecked{Payload: msg})

	return nil
}
