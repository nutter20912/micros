package models

import (
	"context"
	"log"
	mongodb "micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"time"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
