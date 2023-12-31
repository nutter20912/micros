package main

import (
	"context"
	"fmt"
	"log"
	"micros/app/notify/handler"
	"micros/auth"
	_ "micros/broker/natsjs"
	"micros/config"
	"micros/event"
	"micros/otel"
	notifyV1 "micros/proto/notify/v1"
	"micros/wrapper"

	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"github.com/spf13/viper"
	"go-micro.dev/v4/logger"

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

	service.Options().Broker.Connect()

	e := event.New(service.Client())

	notifyV1.RegisterNotifyServiceHandler(
		service.Server(),
		handler.NewNotifyService(service, e))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
