package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var db *redis.Client

func Init() {
	db = newClient()
}

func Get() *redis.Client {
	if db == nil {
		log.Fatal("not init redis client")
	}

	return db
}

func newClient() *redis.Client {
	addr := fmt.Sprintf(
		"%s:%v",
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       viper.GetInt("redis.database"),
		Password: viper.GetString("redis.password")})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("init redis")

	return rdb
}
