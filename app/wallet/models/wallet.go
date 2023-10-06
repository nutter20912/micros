package models

import (
	"context"
	"fmt"
	userV1 "micros/proto/user/v1"
	walletV1 "micros/proto/wallet/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalletEvent struct {
	Id      primitive.ObjectID       `bson:"_id,omitempty"`
	Time    string                   `bson:"event_time,omitempty"`
	Type    walletV1.WalletEventType `bson:"event_type,omitempty"`
	UserId  string                   `bson:"user_id,omitempty"`
	OrderId string                   `bson:"order_id,omitempty"`
	Change  int64                    `bson:"change,omitempty"`
	Memo    string                   `bson:"memo,omitempty"`
}

var (
	Wallet_TABLE = "wallets"
)

func InsertWalletEvent(ctx context.Context, c *mongo.Client, event *userV1.RegisteredEvent) error {
	coll := c.Database("wallet").Collection("wallet_event")

	newWallet := WalletEvent{
		Id:     primitive.NewObjectID(),
		Time:   time.Now().Format(time.RFC3339),
		Type:   walletV1.WalletEventType_WALLET_EVENT_TYPE_SYSTEM,
		UserId: event.UserId,
		Memo:   "init wallet",
	}

	fmt.Println(newWallet)

	if _, err := coll.InsertOne(context.Background(), newWallet); err != nil {
		return err
	}

	return nil
}

func UpdateWallets(ctx context.Context, c *mongo.Client) error {
	coll := c.Database("wallet").Collection("wallet_event")

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$user_id"},
		{Key: "amount", Value: bson.D{{Key: "$sum", Value: "$change"}}},
	}}}

	mergeStage := bson.D{{Key: "$merge", Value: bson.D{
		{Key: "into", Value: Wallet_TABLE},
		{Key: "whenMatched", Value: "replace"},
	}}}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{groupStage, mergeStage})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return err
	}

	for _, result := range results {
		fmt.Println(result)
	}

	return nil
}
