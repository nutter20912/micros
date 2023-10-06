package models

import (
	"context"
	"fmt"
	walletV1 "micros/proto/wallet/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalletEvent struct {
	Id      primitive.ObjectID       `bson:"_id,omitempty"`
	UserId  string                   `bson:"user_id,omitempty"`
	OrderId string                   `bson:"order_id,omitempty"`
	Time    string                   `bson:"event_time,omitempty"`
	Type    walletV1.WalletEventType `bson:"event_type,omitempty"`
	Change  float64                  `bson:"change,omitempty"`
	Memo    string                   `bson:"memo,omitempty"`
}

var (
	Wallet_TABLE = "wallets"
)

func InsertWalletEvent(ctx context.Context, c *mongo.Client, walletEvent *walletV1.WalletEvent) error {
	coll := c.Database("wallet").Collection("wallet_event")

	newWallet := WalletEvent{
		Id:     primitive.NewObjectID(),
		Time:   time.Now().Format(time.RFC3339),
		Type:   walletEvent.Type,
		UserId: walletEvent.UserId,
		Change: walletEvent.Change,
		Memo:   walletEvent.Memo,
	}

	if _, err := coll.InsertOne(context.Background(), newWallet); err != nil {
		return err
	}

	return nil
}

func UpdateWallets(ctx context.Context, c *mongo.Client, userId string) error {
	coll := c.Database("wallet").Collection("wallet_event")

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "user_id", Value: userId},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$user_id"},
		{Key: "amount", Value: bson.D{{Key: "$sum", Value: "$change"}}},
	}}}

	mergeStage := bson.D{{Key: "$merge", Value: bson.D{
		{Key: "into", Value: Wallet_TABLE},
		{Key: "whenMatched", Value: "replace"},
	}}}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, mergeStage})
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
