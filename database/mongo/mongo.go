package mongo

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go-micro.dev/v4/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

func Get() *mongo.Client {
	return db
}

func Init() {
	addr := fmt.Sprintf(
		"mongodb://%s:%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
	)

	options := options.Client().
		ApplyURI(addr).
		SetMaxPoolSize(10).
		SetAuth(options.Credential{
			AuthSource: viper.GetString("db.database"),
			Username:   viper.GetString("db.username"),
			Password:   viper.GetString("db.password"),
		}).
		SetMonitor(newMonitor())

	var err error
	if db, err = mongo.Connect(context.Background(), options); err != nil {
		logger.Fatal(err)
	}

	if err = db.Ping(context.Background(), nil); err != nil {
		logger.Fatal(err)
	}
}
