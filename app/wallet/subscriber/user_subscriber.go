package subscriber

import (
	"context"
	"errors"
	"micros/app/wallet/models"
	userV1 "micros/proto/user/v1"
	walletV1 "micros/proto/wallet/v1"
	"micros/queue"

	"go-micro.dev/v4"
)

type userSubscriber struct {
	Service micro.Service
}

func (s *userSubscriber) initWalletEvent(ctx context.Context, event *userV1.RegisteredEventMessage) error {
	microId, err := queue.MicroId(ctx)
	if err != nil {
		return err
	}

	if err := queue.CheckMsgId(new(models.WalletEvent), microId); err != nil {
		return err
	}

	newWalletEvent := &models.WalletEvent{
		MsgId:  microId,
		Type:   walletV1.WalletEventType_WALLET_EVENT_TYPE_SYSTEM,
		UserId: event.UserId,
		Change: 0,
		Memo:   "init"}

	if err := newWalletEvent.Add(); err != nil {
		return errors.New("add wallet_event error")
	}

	return nil
}
