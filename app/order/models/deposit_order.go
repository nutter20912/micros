package models

import (
	"context"
	"fmt"
	orderV1 "micros/proto/order/v1"
	"time"

	mongodb "micros/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	databaseName                   = "order"
	depositOrderViewCollectionName = "deposit_order_view"
)

type DepositOrder struct {
	Id        string                `json:"id" bson:"id,omitempty"`
	UserId    string                `json:"user_id" bson:"user_id,omitempty"`
	Status    orderV1.DepositStatus `json:"status" bson:"status,omitempty"`
	Amount    float64               `json:"amount" bson:"amount,omitempty"`
	Memo      string                `json:"memo" bson:"memo,omitempty"`
	CreatedAt time.Time             `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time             `json:"updated_at" bson:"updated_at,omitempty"`
}

func (d *DepositOrder) DatabaseName() string {
	return databaseName
}

func (d *DepositOrder) CollectionName() string {
	return depositOrderViewCollectionName
}

func (d *DepositOrder) Get(orderId string) (*DepositOrder, error) {
	coll := mongodb.Get().Database(d.DatabaseName()).Collection(d.CollectionName())

	var order *DepositOrder
	if err := coll.FindOne(context.Background(), bson.M{"id": orderId}).Decode(&order); err != nil {
		return nil, err
	}

	return order, nil
}

func (d *DepositOrder) getAggregatePipeline(orderId string) mongo.Pipeline {
	sortStage := bson.D{{Key: "$sort", Value: bson.D{
		{Key: "time", Value: 1},
	}}}

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "order_id", Value: orderId},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$order_id"},
		{Key: "user_id", Value: bson.D{{Key: "$last", Value: "$user_id"}}},
		{Key: "amount", Value: bson.D{{Key: "$last", Value: "$amount"}}},
		{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$time"}}},
		{Key: "updated_at", Value: bson.D{{Key: "$last", Value: "$time"}}},
		{Key: "status", Value: bson.D{{Key: "$last", Value: "$status"}}},
		{Key: "memo", Value: bson.D{{Key: "$last", Value: "$memo"}}},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "id", Value: "$_id"},
		{Key: "user_id", Value: 1},
		{Key: "amount", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "updated_at", Value: 1},
		{Key: "status", Value: 1},
		{Key: "memo", Value: 1},
	}}}

	return mongo.Pipeline{sortStage, matchStage, groupStage, projectStage}
}

func (d *DepositOrder) Update(ctx context.Context, eventColl *mongo.Collection, orderId string) error {
	mergeStage := bson.D{{Key: "$merge", Value: bson.D{
		{Key: "into", Value: d.CollectionName()},
		{Key: "whenMatched", Value: "replace"},
	}}}

	pipeline := d.getAggregatePipeline(orderId)

	cursor, err := eventColl.Aggregate(ctx, append(pipeline, mergeStage))
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
