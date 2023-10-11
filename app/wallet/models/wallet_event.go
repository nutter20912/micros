package models

import (
	"context"
	mongodb "micros/database/mongo"
	walletV1 "micros/proto/wallet/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	walletViewCollectionName = "wallet_view"
)

type WalletEvent struct {
	Id      primitive.ObjectID       `bson:"_id,omitempty"`
	UserId  string                   `bson:"user_id,omitempty"`
	OrderId string                   `bson:"order_id,omitempty"`
	Time    string                   `bson:"time,omitempty"`
	Type    walletV1.WalletEventType `bson:"type,omitempty"`
	Change  float64                  `bson:"change,omitempty"`
	Memo    string                   `bson:"memo,omitempty"`
}

func (d *WalletEvent) DatabaseName() string {
	return databaseName
}

func (w *WalletEvent) CollectionName() string {
	return walletEventCollectionName
}

func (w *WalletEvent) Add(walletEvent *walletV1.WalletEvent) error {
	coll := mongodb.Get().Database(w.DatabaseName()).Collection(w.CollectionName())

	wallet := WalletEvent{
		Id:      primitive.NewObjectID(),
		Time:    time.Now().Format(time.RFC3339),
		Type:    walletEvent.Type,
		UserId:  walletEvent.UserId,
		OrderId: walletEvent.OrderId,
		Change:  walletEvent.Change,
		Memo:    walletEvent.Memo,
	}

	if _, err := coll.InsertOne(context.Background(), wallet); err != nil {
		return err
	}

	return nil
}
