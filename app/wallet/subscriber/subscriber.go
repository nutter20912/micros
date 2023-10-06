package subscriber

import (
	"context"
	"micros/app/wallet/models"
	"micros/database/mongo"
	userV1 "micros/proto/user/v1"
)

type UserRegisterd struct{}

func (s *UserRegisterd) Handle(ctx context.Context, event *userV1.RegisteredEvent) error {
	client := mongo.Get()

	models.InsertWalletEvent(context.Background(), client, event)

	models.UpdateWallets(context.Background(), client)

	return nil
}
