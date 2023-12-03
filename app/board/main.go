package main

import (
	"context"
	"fmt"
	"log"
	"micros/app/board/handler"
	"micros/auth"
	_ "micros/broker/natsjs"
	"micros/config"
	"micros/database/mysql"
	"micros/otel"
	boardV1 "micros/proto/board/v1"
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
	config.Init("board")
	mysql.Init()
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

	boardV1.RegisterPostServiceHandler(
		service.Server(),
		new(handler.PostService))

	boardV1.RegisterCommentServiceHandler(
		service.Server(),
		new(handler.CommentService))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
