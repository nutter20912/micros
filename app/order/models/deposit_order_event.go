package models

import (
	"context"
	mongodb "micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"time"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	depositOrderEventCollectionName = "deposit_order_event"
)

type DepositOrderEvent struct {
	Id      primitive.ObjectID    `bson:"_id,omitempty"`
	OrderId string                `bson:"order_id,omitempty"`
	UserId  string                `bson:"user_id,omitempty"`
	Status  orderV1.DepositStatus `bson:"status,omitempty"`
	Amount  float64               `bson:"amount,omitempty"`
	Memo    string                `bson:"memo,omitempty"`
	Time    time.Time             `bson:"time,omitempty"`
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
	id := primitive.NewObjectID()

	event := &DepositOrderEvent{
		Id:      id,
		OrderId: orderId,
		UserId:  userId,
		Status:  status,
		Amount:  amount,
		Time:    id.Timestamp(),
	}

	coll := mongodb.Get().Database(d.DatabaseName()).Collection(event.CollectionName())

	if _, err := coll.InsertOne(context.Background(), event); err != nil {
		return nil, err
	}

	return event, nil
}
