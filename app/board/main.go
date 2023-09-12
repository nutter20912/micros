package main

import (
	"fmt"
	"log"
	"micros/app/board/handler"
	"micros/auth"
	"micros/config"
	"micros/database/mysql"
	boardV1 "micros/proto/board/v1"
	"micros/wapper"

	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

	sgrpc "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
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

	re := consul.NewRegistry(registry.Addrs(":8500"))
	a := auth.NewMicroAuth()

	service.Init(
		micro.Server(sgrpc.NewServer(server.Name(appName))),
		micro.Address(fmt.Sprintf(":%s", appPort)),
		micro.Registry(re),
		micro.Auth(a),
		micro.WrapHandler(wapper.NewRequestWrapper()),
		// micro.WrapHandler(wapper.NewAuthWapper(a)),
	)

	boardV1.RegisterPostServiceHandler(
		service.Server(),
		new(handler.BoardService),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
