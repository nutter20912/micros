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

	walletV1 "micros/proto/wallet/v1"
	"micros/wapper"

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
		micro.WrapHandler(wapper.NewRequestWrapper()),
		micro.WrapHandler(wapper.NewAuthWapper(a)),
		micro.WrapSubscriber(wapper.LogSubWapper()))

	walletV1.RegisterWalletServiceHandler(
		service.Server(),
		&handler.WalletService{Service: service})

	micro.RegisterSubscriber(
		"user.*",
		service.Server(),
		&subscriber.UserSubscriber{Service: service})

	micro.RegisterSubscriber(
		"order.*",
		service.Server(),
		&subscriber.OrderSubscriber{Service: service})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
