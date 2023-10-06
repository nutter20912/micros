package subscriber

import (
	"context"
	"micros/app/wallet/models"
	"micros/database/mongo"
	userV1 "micros/proto/user/v1"
	walletV1 "micros/proto/wallet/v1"
)

type UserRegisterd struct{}

func (s *UserRegisterd) Handle(ctx context.Context, event *userV1.RegisteredEvent) error {
	client := mongo.Get()

	walletEvent := walletV1.WalletEvent{
		Type:   walletV1.WalletEventType_WALLET_EVENT_TYPE_SYSTEM,
		UserId: event.UserId,
		Change: 0,
		Memo:   "init",
	}

	models.InsertWalletEvent(context.Background(), client, &walletEvent)
	models.UpdateWallets(context.Background(), client, event.UserId)

	return nil
}
