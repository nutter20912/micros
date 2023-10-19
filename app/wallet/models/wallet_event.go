package models

import (
	"context"
	"log"
	mongodb "micros/database/mongo"
	walletV1 "micros/proto/wallet/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	walletViewCollectionName = "wallet_view"
)

type WalletEvent struct {
	Id      primitive.ObjectID       `json:"id" bson:"_id,omitempty"`
	UserId  string                   `json:"user_id" bson:"user_id,omitempty"`
	OrderId string                   `json:"order_id" bson:"order_id,omitempty"`
	Time    time.Time                `json:"time" bson:"time,omitempty"`
	Type    walletV1.WalletEventType `json:"type" bson:"type,omitempty"`
	Change  float64                  `json:"change" bson:"change,omitempty"`
	Memo    string                   `json:"memo" bson:"memo,omitempty"`
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
		Time:    time.Now(),
		Type:    walletEvent.Type,
		UserId:  walletEvent.UserId,
		OrderId: walletEvent.OrderId,
		Change:  walletEvent.Change,
		Memo:    walletEvent.Memo,
	}

	if _, err := coll.InsertOne(context.Background(), wallet); err != nil {
		return err
	}

	if err := new(Wallet).Update(context.Background(), coll, wallet.UserId); err != nil {
		log.Printf("[wallet err]: %v", err.Error())
	}

	return nil
}

func (w *WalletEvent) Get(userId string, page *int64, limit *int64) ([]*WalletEvent, *mongodb.Paginatior, error) {
	coll := mongodb.Get().Database(w.DatabaseName()).Collection(w.CollectionName())

	var events []*WalletEvent

	paginatior, err := mongodb.NewPagination(coll).
		Where(bson.M{"user_id": userId}).
		Desc("_id").
		Page(page).
		Limit(limit).
		Find(context.Background(), &events)
	if err != nil {
		return nil, nil, err
	}

	return events, paginatior, nil
}
