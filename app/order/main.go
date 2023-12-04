package main

import (
	"context"
	"fmt"
	"log"
	"micros/app/order/handler"
	"micros/app/order/subscriber"
	"micros/auth"
	_ "micros/broker/natsjs"
	"micros/config"
	"micros/database/mongo"
	"micros/event"
	"micros/otel"
	orderV1 "micros/proto/order/v1"
	"micros/wrapper"

	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"github.com/spf13/viper"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
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
	otelShutdown, err := otel.SetupGlobalOTelSDK(context.Background(), appName, "0.1.0")
	if err != nil {
		logger.Error(err)
	}

	a := auth.NewMicroAuth()

	service.Init(
		micro.Server(sgrpc.NewServer(server.Name(appName))),
		micro.Address(fmt.Sprintf(":%s", appPort)),
		micro.Auth(a),
		micro.BeforeStop(func() error {
			return otelShutdown(context.Background())
		}),

		micro.WrapClient(
			opentelemetry.NewClientWrapper(),
			wrapper.NewClientWrapper),

		micro.WrapHandler(
			opentelemetry.NewHandlerWrapper(),
			wrapper.NewRequestWrapper(),
			wrapper.NewAuthWrapper(a)),

		micro.WrapSubscriber(
			opentelemetry.NewSubscriberWrapper(),
			wrapper.LogSubWrapper()))

	e := event.New(service.Client())

	orderV1.RegisterOrderServiceHandler(
		service.Server(),
		handler.NewOrderService(service, e))

	subscriber.Register(service, e)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
