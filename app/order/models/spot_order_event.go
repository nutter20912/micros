package models

import (
	"context"
	"log"
	mongodb "micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"time"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SpotOrderEvent struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderId string             `json:"order_id" bson:"order_id,omitempty"`
	UserId  string             `json:"user_id" bson:"user_id,omitempty"`
	Time    time.Time          `json:"time" bson:"time,omitempty"`

	Symbol   string             `json:"symbol" bson:"symbol,omitempty"`
	Quantity float64            `json:"quantity" bson:"quantity,omitempty"`
	Side     orderV1.OrderSide  `json:"side" bson:"side,omitempty"`
	Type     orderV1.OrderType  `json:"type" bson:"type,omitempty"`
	Price    float64            `json:"price" bson:"price,omitempty"`
	Status   orderV1.SpotStatus `json:"status" bson:"status,omitempty"`
	Memo     string             `json:"memo" bson:"memo,omitempty"`

	MsgId string `json:"msg_id" bson:"msg_id,omitempty"`
}

func (d *SpotOrderEvent) DatabaseName() string {
	return databaseName
}

func (d *SpotOrderEvent) CollectionName() string {
	return "spot_order_event"
}

func (e *SpotOrderEvent) Create() error {
	e.OrderId = ulid.Make().String()

	if err := e.Add(); err != nil {
		return err
	}

	return nil
}

func (e *SpotOrderEvent) Add() error {
	coll := mongodb.Get().Database(e.DatabaseName()).Collection(e.CollectionName())

	id := primitive.NewObjectID()
	e.Id = id
	e.Time = id.Timestamp()

	if _, err := coll.InsertOne(context.Background(), e); err != nil {
		return err
	}

	if err := new(SpotOrder).Update(context.Background(), coll, e.OrderId); err != nil {
		log.Printf("[deposit order err]: %v", err.Error())
	}

	return nil
}

func (e *SpotOrderEvent) Last() error {
	coll := mongodb.Get().Database(e.DatabaseName()).Collection(e.CollectionName())

	opts := options.FindOne().SetSort(bson.M{"time": -1})
	if err := coll.FindOne(context.Background(), bson.M{"order_id": e.OrderId}, opts).Decode(e); err != nil {
		return err
	}

	return nil
}

func (e *SpotOrderEvent) Exist(microId string) (bool, error) {
	coll := mongodb.Get().Database(e.DatabaseName()).Collection(e.CollectionName())

	count, err := coll.CountDocuments(context.Background(), bson.M{"msg_id": microId})
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (e *SpotOrderEvent) Get() (res *SpotOrderEvent, err error) {
	coll := mongodb.Get().Database(e.DatabaseName()).Collection(e.CollectionName())

	filter := bson.M{"order_id": e.OrderId, "user_id": e.UserId}

	if e.Id != primitive.NilObjectID {
		filter["_id"] = bson.M{"$gt": e.Id}
	}

	if err := coll.FindOne(context.Background(), filter).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (e *SpotOrderEvent) Count() (int64, error) {
	coll := mongodb.Get().Database(e.DatabaseName()).Collection(e.CollectionName())

	filter := bson.M{}

	if e.OrderId != "" {
		filter["order_id"] = e.OrderId
	}

	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return count, err
}
