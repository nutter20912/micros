package subscriber

import (
	"context"
	"errors"
	"micros/app/wallet/models"
	userV1 "micros/proto/user/v1"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
)

type UserSubscriber struct {
	Service micro.Service
}

func (s *UserSubscriber) Registered(ctx context.Context, event *userV1.RegisteredEventMessage) error {
	walletEvent := walletV1.WalletEvent{
		Type:   walletV1.WalletEventType_WALLET_EVENT_TYPE_SYSTEM,
		UserId: event.UserId,
		Change: 0,
		Memo:   "init",
	}

	if err := new(models.WalletEvent).Add(&walletEvent); err != nil {
		return errors.New("add wallet_event error")
	}

	return nil
}
