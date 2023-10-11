package models

import (
	"context"
	"fmt"
	orderV1 "micros/proto/order/v1"
	"time"

	mongodb "micros/database/mongo"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	databaseName                    = "order"
	depositOrderEventCollectionName = "deposit_order_event"
	depositOrderViewCollectionName  = "deposit_order_view"
)

type Model interface {
	CollectionName() string
}

type DepositOrderEvent struct {
	Id      primitive.ObjectID    `bson:"_id,omitempty"`
	OrderId string                `bson:"order_id,omitempty"`
	UserId  string                `bson:"user_id,omitempty"`
	Status  orderV1.DepositStatus `bson:"status,omitempty"`
	Amount  float64               `bson:"amount,omitempty"`
	Memo    string                `bson:"memo,omitempty"`
	Time    time.Time             `bson:"time,omitempty"`
}

func (d DepositOrderEvent) CollectionName() string {
	return depositOrderEventCollectionName
}

type DepositOrder struct {
	Id     string                `bson:"id,omitempty"`
	UserId string                `bson:"user_id,omitempty"`
	Status orderV1.DepositStatus `bson:"status,omitempty"`
	Amount float64               `bson:"amount,omitempty"`
	Memo   string                `bson:"memo,omitempty"`
}

func (d DepositOrder) CollectionName() string {
	return depositOrderViewCollectionName
}

func CreateDepositOrderEvent(userId string, amount float64) (event *DepositOrderEvent, err error) {
	orderId := ulid.Make().String()

	if event, err = InertDepositOrderEvent(orderId, userId, amount, orderV1.DepositStatus_DEPOSIT_STATUS_PROCESSING); err != nil {
		return nil, err
	}

	return event, nil
}

func InertDepositOrderEvent(
	orderId string,
	userId string,
	amount float64,
	status orderV1.DepositStatus,
) (*DepositOrderEvent, error) {
	id := primitive.NewObjectID()

	event := &DepositOrderEvent{
		Id:      id,
		OrderId: orderId,
		UserId:  userId,
		Status:  status,
		Amount:  amount,
		Time:    id.Timestamp(),
	}

	coll := mongodb.Get().Database(databaseName).Collection(event.CollectionName())

	if _, err := coll.InsertOne(context.Background(), event); err != nil {
		return nil, err
	}

	return event, nil
}

func GetDepositOrder(ctx context.Context, orderId string) (*DepositOrder, error) {
	d := new(DepositOrder)

	coll := mongodb.Get().Database(databaseName).Collection(d.CollectionName())

	if err := coll.FindOne(ctx, bson.M{"id": orderId}).Decode(&d); err != nil {
		return nil, err
	}

	return d, nil
}

func UpdateOrder(ctx context.Context, eventColl *mongo.Collection, orderId string) error {
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
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "id", Value: "$_id"},
		{Key: "user_id", Value: 1},
		{Key: "amount", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "updated_at", Value: 1},
		{Key: "status", Value: 1},
	}}}

	mergeStage := bson.D{{Key: "$merge", Value: bson.D{
		{Key: "into", Value: new(DepositOrder).CollectionName()},
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
