package models

import (
	"context"
	"fmt"
	"time"

	mongodb "micros/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	databaseName              = "wallet"
	walletEventCollectionName = "wallet_event"
)

type Wallet struct {
	UserId    string    `json:"user_id" bson:"user_id,omitempty"`
	Amount    float64   `json:"amount" bson:"amount,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

func (d *Wallet) DatabaseName() string {
	return databaseName
}

func (w *Wallet) CollectionName() string {
	return walletViewCollectionName
}

func (w *Wallet) Get(ctx context.Context, userId string) (*Wallet, error) {
	coll := mongodb.Get().Database(w.DatabaseName()).Collection(w.CollectionName())

	var wallet Wallet
	if err := coll.FindOne(ctx, bson.M{"_id": userId}).Decode(&wallet); err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (w *Wallet) Update(ctx context.Context, eventColl *mongo.Collection, userId string) error {
	sortStage := bson.D{{Key: "$sort", Value: bson.D{
		{Key: "time", Value: 1},
	}}}

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "user_id", Value: userId},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$user_id"},
		{Key: "amount", Value: bson.D{{Key: "$sum", Value: "$change"}}},
		{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$time"}}},
		{Key: "updated_at", Value: bson.D{{Key: "$last", Value: "$time"}}},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "user_id", Value: "$_id.user_id"},
		{Key: "amount", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "updated_at", Value: 1},
	}}}

	mergeStage := bson.D{{Key: "$merge", Value: bson.D{
		{Key: "into", Value: w.CollectionName()},
		{Key: "whenMatched", Value: "replace"},
	}}}

	cursor, err := eventColl.Aggregate(ctx, mongo.Pipeline{sortStage, matchStage, groupStage, projectStage, mergeStage})
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
