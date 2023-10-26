package subscriber

import (
	"context"
	"micros/app/wallet/event"
	"micros/app/wallet/models"
	baseEvent "micros/event"
	orderV1 "micros/proto/order/v1"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
)

type OrderSubscriber struct {
	Service micro.Service
}

func (o *OrderSubscriber) DepositCreated(ctx context.Context, e *orderV1.DepositOrderEvent) error {
	microId, err := baseEvent.MicroId(ctx)
	if err != nil {
		return err
	}

	if err := validate(ctx, microId); err != nil {
		return err
	}

	getMessage := func(e *orderV1.DepositOrderEvent) *walletV1.TransactionEventMessage {
		msg := &walletV1.TransactionEventMessage{
			UserId:  e.UserId,
			OrderId: e.OrderId,
			Type:    walletV1.WalletEventType_WALLET_EVENT_TYPE_DEPOSIT,
			Success: true}

		if _, err := new(models.Wallet).Get(e.UserId); err != nil {
			msg.Memo = "wallet not found"
			msg.Success = false
			return msg
		}

		newWalletEvent := &models.WalletEvent{
			MsgId:   microId,
			UserId:  e.UserId,
			OrderId: e.OrderId,
			Change:  e.Amount,
			Type:    walletV1.WalletEventType_WALLET_EVENT_TYPE_DEPOSIT}

		if err := new(models.WalletEvent).Add(newWalletEvent); err != nil {
			msg.Memo = "add wallet_event error"
			msg.Success = false
			return msg
		}

		return msg
	}

	event.TransactionEvent{Client: o.Service.Client()}.Dispatch(getMessage(e))

	return nil
}

func (o *OrderSubscriber) OrderCheck(ctx context.Context, e *orderV1.OrderCheckEventMessage) error {
	msg := orderV1.CheckCallbackMessage{
		Type: orderV1.CheckCallbackType_CHECK_CALLBACK_TYPE_WALLET,
	}

	wallet, err := new(models.Wallet).Get(e.UserId)
	if err != nil {
		msg.Success = false
		msg.Msg = err.Error()
	}

	if wallet.Amount == 0 {
		msg.Success = false
		msg.Msg = "balance not enough"
	} else {
		msg.Success = true
		msg.Data = &orderV1.CheckCallbackMessage_Wallet{
			Wallet: &orderV1.CheckWalletData{
				Balance: wallet.Amount,
			},
		}
	}

	event.CheckCallbackEvent{Client: o.Service.Client()}.Dispatch(e.CallbackSubject, &msg)

	return nil
}
