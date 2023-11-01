package main

import (
	"fmt"
	"log"
	"micros/app/wallet/handler"
	"micros/app/wallet/subscriber"
	"micros/auth"
	_ "micros/broker/natsjs"
	"micros/config"
	"micros/database/mongo"
	"micros/event"

	walletV1 "micros/proto/wallet/v1"
	"micros/wrapper"

	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

func init() {
	config.Init("wallet")
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
		micro.WrapHandler(wrapper.NewRequestWrapper()),
		micro.WrapHandler(wrapper.NewAuthWrapper(a)),
		micro.WrapSubscriber(wrapper.LogSubWrapper()))

	e := event.New(service.Client())

	walletV1.RegisterWalletServiceHandler(
		service.Server(),
		handler.NewWalletService(service, e))

	subscriber.Register(service, e)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
