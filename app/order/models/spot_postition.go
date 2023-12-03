package models

import (
	"context"
	"errors"
	"fmt"
	mongodb "micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SpotPosition struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    string             `json:"user_id" bson:"user_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`

	Symbol   string            `json:"symbol" bson:"symbol,omitempty"`
	Side     orderV1.OrderSide `json:"side" bson:"side,omitempty"`
	Quantity float64           `json:"quantity" bson:"quantity,omitempty"`

	OrderId      string  `json:"order_id" bson:"order_id,omitempty"`
	Price        float64 `json:"price" bson:"price,omitempty"`
	Fee          float64 `json:"fee" bson:"fee,omitempty"`
	OpenQuantity float64 `json:"open_quantity" bson:"open_quantity,omitempty"`
}

func (s *SpotPosition) DatabaseName() string {
	return databaseName
}

func (s *SpotPosition) CollectionName() string {
	return "spot_position"
}

func (s *SpotPosition) Get(ctx context.Context, userId string, symbol string, page *int64, limit *int64) ([]*SpotPosition, *mongodb.Paginator, error) {
	coll := mongodb.Get().Database(s.DatabaseName()).Collection(s.CollectionName())

	filter := bson.M{
		"symbol":        symbol,
		"user_id":       userId,
		"open_quantity": bson.M{"$gt": 0},
	}

	var data []*SpotPosition

	paginator, err := mongodb.NewPagination(coll).
		Where(filter).
		Desc("_id").
		Page(page).
		Limit(limit).
		Find(ctx, &data)
	if err != nil {
		return nil, nil, err
	}

	return data, paginator, nil
}

func (s *SpotPosition) GetList(ctx context.Context, userId string, symbol string) ([]*SpotPosition, error) {
	coll := mongodb.Get().Database(s.DatabaseName()).Collection(s.CollectionName())

	filter := bson.M{
		"symbol":        symbol,
		"user_id":       userId,
		"open_quantity": bson.M{"$gt": 0},
	}

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})

	var data []*SpotPosition
	cur, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &data); err != nil {
		return nil, err
	}

	bson.MarshalExtJSON(data, false, false)

	return data, nil
}

func (s *SpotPosition) Upsert(ctx context.Context) error {
	coll := mongodb.Get().Database(s.DatabaseName()).Collection(s.CollectionName())

	filter := bson.M{
		"user_id":       s.UserId,
		"symbol":        s.Symbol,
		"open_quantity": bson.M{"$gt": 0},
	}

	if s.Side == orderV1.OrderSide_ORDER_SIDE_BUY {
		filter["side"] = orderV1.OrderSide_ORDER_SIDE_SELL
	} else {
		filter["side"] = orderV1.OrderSide_ORDER_SIDE_BUY
	}

	// 撈庫存
	var positions []*SpotPosition
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}

	if err = cur.All(ctx, &positions); err != nil {
		return err
	}

	//建倉
	if errors.Is(mongo.ErrNoDocuments, err) {
		if err := s.Add(ctx); err != nil {
			return err
		}

		return nil
	}

	//平倉
	s.close(ctx, positions)

	if s.Quantity > 0 {
		s.Add(ctx)
		return nil
	}

	return err
}

func (s *SpotPosition) close(ctx context.Context, positions []*SpotPosition) {
	for _, position := range positions {
		switch {
		//完全平倉
		case position.OpenQuantity < s.Quantity:
			s.Quantity -= position.OpenQuantity

			position.CloseEvent(ctx, s, position.OpenQuantity)
			position.UpdateOpenQuantity(ctx, 0)

		//部分平倉
		case position.OpenQuantity > s.Quantity:
			position.CloseEvent(ctx, s, s.Quantity)
			position.UpdateOpenQuantity(ctx, position.OpenQuantity-s.Quantity)

			s.Quantity = 0
			return

		//相等
		default:
			position.CloseEvent(ctx, s, s.Quantity)
			position.UpdateOpenQuantity(ctx, 0)
			s.Quantity = 0
			return
		}
	}
}

func (s *SpotPosition) Add(ctx context.Context) error {
	coll := mongodb.Get().Database(s.DatabaseName()).Collection(s.CollectionName())

	id := primitive.NewObjectID()
	s.Id = id
	s.CreatedAt = id.Timestamp()
	s.OpenQuantity = s.Quantity

	if _, err := coll.InsertOne(ctx, s); err != nil {
		return err
	}
	return nil
}

func (s *SpotPosition) UpdateOpenQuantity(ctx context.Context, openQuantity float64) error {
	coll := mongodb.Get().Database(s.DatabaseName()).Collection(s.CollectionName())

	filter := bson.M{"_id": s.Id}
	update := bson.M{"$set": bson.M{
		"open_quantity": openQuantity,
		"updated_at":    time.Now(),
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&s); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *SpotPosition) CloseEvent(ctx context.Context, incoming *SpotPosition, closeQty float64) {
	spc := SpotPositionClosed{
		UserId:   s.UserId,
		Symbol:   s.Symbol,
		Side:     s.Side,
		Quantity: closeQty,

		OpenOrderId: s.OrderId,
		OpenPrice:   s.Price,
		OpenFee:     s.Fee,

		CloseOrderId: incoming.OrderId,
		ClosePrice:   incoming.Price,
		CloseFee:     incoming.Fee,
	}

	if err := spc.Add(ctx); err != nil {
		fmt.Println(err)
	}
}
