package main

import (
	"fmt"
	"log"
	"micros/app/order/handler"
	"micros/app/order/subscriber"
	"micros/auth"
	"micros/config"
	"micros/database/mongo"
	orderV1 "micros/proto/order/v1"
	"micros/wapper"

	_ "github.com/go-micro/plugins/v4/broker/nats"
	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

func init() {
	config.Init("order")
	mongo.Init()
}

func main() {
	appName := viper.GetString("app.name")
	appPort := viper.GetString("app.port")
	service := micro.NewService()

	a := auth.NewMicroAuth()

	service.Init(
		micro.Server(sgrpc.NewServer(server.Name(appName))),
		micro.Address(fmt.Sprintf(":%s", appPort)),
		micro.Auth(a),
		micro.WrapHandler(wapper.NewRequestWrapper()),
		micro.WrapHandler(wapper.NewAuthWapper(a)),
		micro.WrapSubscriber(wapper.LogSubWapper()))

	orderV1.RegisterOrderServiceHandler(
		service.Server(),
		&handler.OrderService{Service: service})

	micro.RegisterSubscriber(
		viper.GetString("topic.wallet.transaction"),
		service.Server(),
		&subscriber.AddOrderEvent{})

	//go models.OrderWatcher(mongo.Get())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
