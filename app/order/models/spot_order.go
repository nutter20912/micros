package models

import (
	"context"
	"fmt"
	mongodb "micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpotOrder struct {
	Id        string    `json:"id" bson:"id,omitempty"`
	UserId    string    `json:"user_id" bson:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`

	Symbol   string             `json:"symbol" bson:"symbol,omitempty"`
	Quantity float64            `json:"quantity" bson:"quantity,omitempty"`
	Side     orderV1.OrderSide  `json:"side" bson:"side,omitempty"`
	Type     orderV1.OrderType  `json:"type" bson:"type,omitempty"`
	Price    float64            `json:"price" bson:"price,omitempty"`
	Status   orderV1.SpotStatus `json:"status" bson:"status,omitempty"`
	Memo     string             `json:"memo" bson:"memo,omitempty"`
}

func (s *SpotOrder) DatabaseName() string {
	return databaseName
}

func (s *SpotOrder) CollectionName() string {
	return "spot_order_view"
}

func (e *SpotOrder) Get(
	ctx context.Context,
	page *int64,
	limit *int64,
	filterOptions ...mongodb.FilterOption,
) ([]*SpotOrder, *mongodb.Paginator, error) {
	coll := mongodb.Get().Database(e.DatabaseName()).Collection(e.CollectionName())

	var events []*SpotOrder
	filter := bson.M{}
	for _, o := range filterOptions {
		o(filter)
	}

	paginator, err := mongodb.NewPagination(coll).
		Where(filter).
		Desc("_id").
		Page(page).
		Limit(limit).
		Find(ctx, &events)
	if err != nil {
		return nil, nil, err
	}

	return events, paginator, nil
}

func (d *SpotOrder) getAggregatePipeline(orderId string) mongo.Pipeline {
	sortStage := bson.D{{Key: "$sort", Value: bson.D{
		{Key: "time", Value: 1},
	}}}

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "order_id", Value: orderId},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$order_id"},
		{Key: "user_id", Value: bson.D{{Key: "$last", Value: "$user_id"}}},
		{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$time"}}},
		{Key: "updated_at", Value: bson.D{{Key: "$last", Value: "$time"}}},

		{Key: "symbol", Value: bson.D{{Key: "$last", Value: "$symbol"}}},
		{Key: "quantity", Value: bson.D{{Key: "$last", Value: "$quantity"}}},
		{Key: "side", Value: bson.D{{Key: "$last", Value: "$side"}}},
		{Key: "type", Value: bson.D{{Key: "$last", Value: "$type"}}},
		{Key: "price", Value: bson.D{{Key: "$last", Value: "$price"}}},
		{Key: "status", Value: bson.D{{Key: "$last", Value: "$status"}}},
		{Key: "memo", Value: bson.D{{Key: "$last", Value: "$memo"}}},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "id", Value: "$_id"},
		{Key: "user_id", Value: 1},
		{Key: "created_at", Value: 1},
		{Key: "updated_at", Value: 1},

		{Key: "symbol", Value: 1},
		{Key: "quantity", Value: 1},
		{Key: "side", Value: 1},
		{Key: "type", Value: 1},
		{Key: "price", Value: 1},
		{Key: "status", Value: 1},
		{Key: "memo", Value: 1},
	}}}

	return mongo.Pipeline{sortStage, matchStage, groupStage, projectStage}
}

func (s *SpotOrder) Update(ctx context.Context, eventColl *mongo.Collection, orderId string) error {
	mergeStage := bson.D{{Key: "$merge", Value: bson.D{
		{Key: "into", Value: s.CollectionName()},
		{Key: "whenMatched", Value: "replace"},
	}}}

	pipeline := s.getAggregatePipeline(orderId)

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
