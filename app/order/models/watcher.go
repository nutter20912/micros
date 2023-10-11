package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderWatcher(client *mongo.Client) {
	database := client.Database(databaseName)
	eventsCollection := database.Collection(new(DepositOrderEvent).CollectionName())
	pipeline := []bson.M{{"$match": bson.M{"operationType": "insert"}}}

	changeStream, err := eventsCollection.Watch(context.TODO(), pipeline)
	if err != nil {
		log.Fatalln(err)
	}

	defer changeStream.Close(context.TODO())

	for changeStream.Next(context.Background()) {
		var event struct {
			Doc DepositOrderEvent `bson:"fullDocument"`
		}

		if err := changeStream.Decode(&event); err != nil {
			fmt.Println("Error decoding event:", err)
			continue
		}

		if err := new(DepositOrder).Update(context.Background(), eventsCollection, event.Doc.OrderId); err != nil {
			fmt.Println(err)
		}
	}

	if err := changeStream.Err(); err != nil {
		fmt.Println("Change stream error:", err)
	}

}
