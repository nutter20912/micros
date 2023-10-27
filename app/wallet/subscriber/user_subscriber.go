package subscriber

import (
	"context"
	"errors"
	"micros/app/wallet/models"
	baseEvent "micros/event"
	userV1 "micros/proto/user/v1"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
)

type userSubscriber struct {
	Service micro.Service
}

func (s *userSubscriber) registered(ctx context.Context, event *userV1.RegisteredEventMessage) error {
	microId, err := baseEvent.MicroId(ctx)
	if err != nil {
		return err
	}

	if err := validate(ctx, microId); err != nil {
		return err
	}

	newWalletEvent := &models.WalletEvent{
		MsgId:  microId,
		Type:   walletV1.WalletEventType_WALLET_EVENT_TYPE_SYSTEM,
		UserId: event.UserId,
		Change: 0,
		Memo:   "init"}

	if err := new(models.WalletEvent).Add(newWalletEvent); err != nil {
		return errors.New("add wallet_event error")
	}

	return nil
}
