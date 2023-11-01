package main

import (
	"fmt"
	"log"
	"micros/app/board/handler"
	"micros/auth"
	"micros/config"
	"micros/database/mysql"
	boardV1 "micros/proto/board/v1"
	"micros/wrapper"

	_ "github.com/go-micro/plugins/v4/broker/nats"
	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

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

	a := auth.NewMicroAuth()

	service.Init(
		micro.Server(sgrpc.NewServer(server.Name(appName))),
		micro.Address(fmt.Sprintf(":%s", appPort)),
		micro.Auth(a),
		micro.WrapHandler(wrapper.NewRequestWrapper()),
		micro.WrapHandler(wrapper.NewAuthWrapper(a)),
	)

	boardV1.RegisterPostServiceHandler(
		service.Server(),
		new(handler.PostService),
	)
	boardV1.RegisterCommentServiceHandler(
		service.Server(),
		new(handler.CommentService),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
