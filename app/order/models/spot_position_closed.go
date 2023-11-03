package models

import (
	"context"
	mongodb "micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 平倉明細
type SpotPositionClosed struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    string             `json:"user_id" bson:"user_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`

	Symbol   string            `json:"symbol" bson:"symbol,omitempty"`
	Side     orderV1.OrderSide `json:"side" bson:"side,omitempty"`
	Quantity float64           `json:"quantity" bson:"quantity,omitempty"`

	OpenOrderId string  `json:"open_order_id" bson:"open_order_id,omitempty"`
	OpenPrice   float64 `json:"open_price" bson:"open_price,omitempty"`
	OpenFee     float64 `json:"open_fee" bson:"open_fee,omitempty"`

	CloseOrderId string  `json:"close_order_id" bson:"close_order_id,omitempty"`
	ClosePrice   float64 `json:"close_price" bson:"close_price,omitempty"`
	CloseFee     float64 `json:"close_fee" bson:"close_fee,omitempty"`
}

func (s *SpotPositionClosed) DatabaseName() string {
	return databaseName
}

func (s *SpotPositionClosed) CollectionName() string {
	return "spot_position_closed"
}

func (s *SpotPositionClosed) Add() error {
	coll := mongodb.Get().Database(s.DatabaseName()).Collection(s.CollectionName())

	id := primitive.NewObjectID()
	s.Id = id
	s.CreatedAt = id.Timestamp()

	if _, err := coll.InsertOne(context.Background(), s); err != nil {
		return err
	}
	return nil
}

func (s *SpotPositionClosed) GetList(userId string) ([]*SpotPositionClosed, error) {
	coll := mongodb.Get().Database(s.DatabaseName()).Collection(s.CollectionName())

	filter := bson.M{
		"user_id": userId,
		//"symbol":  symbol,
	}

	var data []*SpotPositionClosed
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	cur, err := coll.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	bson.MarshalExtJSON(data, false, false)

	return data, nil
}
