package main

import (
	"context"
	"fmt"
	"log"
	"micros/app/market/handler"
	"micros/app/market/subscriber"
	"micros/app/market/tasks"
	"micros/auth"
	"micros/config"
	"micros/database/redis"
	"micros/event"
	"micros/otel"
	marketV1 "micros/proto/market/v1"
	"micros/wrapper"

	_ "micros/broker/natsjs"

	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"github.com/spf13/viper"
	"go-micro.dev/v4/logger"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

func init() {
	config.Init("market")
	redis.Init()
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

	marketV1.RegisterMarketServiceHandler(
		service.Server(),
		handler.NewMarketService(service, e))

	subscriber.Register(service, e)

	go tasks.ExecuteAll(service)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
