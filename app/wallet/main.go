package main

import (
	"context"
	"fmt"
	"log"
	"micros/app/wallet/handler"
	"micros/app/wallet/subscriber"
	"micros/auth"
	_ "micros/broker/natsjs"
	"micros/config"
	"micros/database/mongo"
	"micros/event"
	"micros/otel"

	walletV1 "micros/proto/wallet/v1"
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
	config.Init("wallet")
	mongo.Init()
}

func main() {
	appName := viper.GetString("app.name")
	appPort := viper.GetString("app.port")
	service := micro.NewService()
	otelShutdown, err := otel.SetupOTelSDK(context.Background(), appName, "0.1.0")
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

		micro.WrapClient(opentelemetry.NewClientWrapper()),
		micro.WrapClient(wrapper.NewClientWrapper),

		micro.WrapHandler(opentelemetry.NewHandlerWrapper()),
		micro.WrapHandler(wrapper.NewRequestWrapper()),
		micro.WrapHandler(wrapper.NewAuthWrapper(a)),

		micro.WrapSubscriber(opentelemetry.NewSubscriberWrapper()),
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
