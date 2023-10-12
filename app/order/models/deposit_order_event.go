package models

import (
	"context"
	mongodb "micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"time"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	depositOrderEventCollectionName = "deposit_order_event"
)

type DepositOrderEvent struct {
	Id      primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	OrderId string                `json:"order_id" bson:"order_id,omitempty"`
	UserId  string                `json:"user_id" bson:"user_id,omitempty"`
	Status  orderV1.DepositStatus `json:"status" bson:"status,omitempty"`
	Amount  float64               `json:"amount" bson:"amount,omitempty"`
	Memo    string                `json:"memo" bson:"memo,omitempty"`
	Time    time.Time             `json:"time" bson:"time,omitempty"`
}

func (d *DepositOrderEvent) DatabaseName() string {
	return databaseName
}

func (d *DepositOrderEvent) CollectionName() string {
	return depositOrderEventCollectionName
}

func (d *DepositOrderEvent) Create(userId string, amount float64) (event *DepositOrderEvent, err error) {
	orderId := ulid.Make().String()

	if event, err = d.Add(orderId, userId, amount, orderV1.DepositStatus_DEPOSIT_STATUS_PROCESSING); err != nil {
		return nil, err
	}

	return event, nil
}

func (d *DepositOrderEvent) Add(
	orderId string,
	userId string,
	amount float64,
	status orderV1.DepositStatus,
) (*DepositOrderEvent, error) {
	coll := mongodb.Get().Database(d.DatabaseName()).Collection(d.CollectionName())

	id := primitive.NewObjectID()

	event := &DepositOrderEvent{
		Id:      id,
		OrderId: orderId,
		UserId:  userId,
		Status:  status,
		Amount:  amount,
		Time:    id.Timestamp(),
	}

	if _, err := coll.InsertOne(context.Background(), event); err != nil {
		return nil, err
	}

	return event, nil
}

func (d *DepositOrderEvent) Get(orderId string) ([]*DepositOrderEvent, error) {
	coll := mongodb.Get().Database(d.DatabaseName()).Collection(d.CollectionName())

	var events []*DepositOrderEvent
	cur, err := coll.Find(context.Background(), bson.M{"order_id": orderId})
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.Background(), &events); err != nil {
		return nil, err
	}

	bson.MarshalExtJSON(events, false, false)

	return events, nil
}
