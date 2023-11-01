package main

import (
	"fmt"
	"log"
	"micros/app/user/handler"
	"micros/auth"
	_ "micros/broker/natsjs"
	"micros/config"
	"micros/database/mysql"
	userV1 "micros/proto/user/v1"
	"micros/wrapper"

	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

func init() {
	config.Init("user")
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
		micro.WrapHandler(wrapper.NewAuthWrapper(a)))

	userV1.RegisterUserServiceHandler(
		service.Server(),
		handler.NewUserService(service, mysql.Get()))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
