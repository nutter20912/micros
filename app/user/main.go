package main

import (
	"fmt"
	"log"
	"micros/app/user/handler"
	"micros/auth"
	"micros/config"
	"micros/database/mysql"
	userV1 "micros/proto/user/v1"
	"micros/wapper"

	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
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

	re := consul.NewRegistry(registry.Addrs(":8500"))
	a := auth.NewMicroAuth()

	service.Init(
		micro.Server(sgrpc.NewServer(server.Name(appName))),
		micro.Address(fmt.Sprintf(":%s", appPort)),
		micro.Registry(re),
		micro.Auth(a),
		micro.WrapHandler(wapper.NewRequestWrapper()),
		micro.WrapHandler(wapper.NewAuthWapper(a)),
	)

	userV1.RegisterUserServiceHandler(
		service.Server(),
		new(handler.UserService),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
