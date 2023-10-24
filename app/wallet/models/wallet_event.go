package models

import (
	"context"
	"log"
	mongodb "micros/database/mongo"
	walletV1 "micros/proto/wallet/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	MsgId string `json:"msg_id" bson:"msg_id,omitempty"`
}

func (d *WalletEvent) DatabaseName() string {
	return databaseName
}

func (w *WalletEvent) CollectionName() string {
	return walletEventCollectionName
}

func (w *WalletEvent) Add(we *WalletEvent) error {
	coll := mongodb.Get().Database(w.DatabaseName()).Collection(w.CollectionName())

	we.Id = primitive.NewObjectID()
	we.Time = time.Now()

	if _, err := coll.InsertOne(context.Background(), we); err != nil {
		return err
	}

	if err := new(Wallet).Update(context.Background(), coll, we.UserId); err != nil {
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

func (w *WalletEvent) GetEvents(userId string, eventCursor *string) ([]*WalletEvent, error) {
	coll := mongodb.Get().Database(w.DatabaseName()).Collection(w.CollectionName())

	var events []*WalletEvent

	filter := bson.M{"user_id": userId}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})

	if eventCursor != nil {
		objID, err := primitive.ObjectIDFromHex(*eventCursor)
		if err != nil {
			return nil, err
		}

		filter["_id"] = bson.M{"$gt": objID}
	} else {
		opts.SetLimit(1)
	}

	cur, err := coll.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(context.Background(), &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (w *WalletEvent) Exist(msgId string) (bool, error) {
	coll := mongodb.Get().Database(w.DatabaseName()).Collection(w.CollectionName())

	count, err := coll.CountDocuments(context.Background(), bson.M{"msg_id": msgId})
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
