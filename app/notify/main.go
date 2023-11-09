package main

import (
	"fmt"
	"log"
	"micros/app/notify/handler"
	"micros/auth"
	_ "micros/broker/natsjs"
	"micros/config"
	"micros/event"
	notifyV1 "micros/proto/notify/v1"
	"micros/wrapper"

	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

func init() {
	config.Init("notify")
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

	service.Options().Broker.Connect()

	e := event.New(service.Client())

	notifyV1.RegisterNotifyServiceHandler(
		service.Server(),
		handler.NewNotifyService(service, e))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
