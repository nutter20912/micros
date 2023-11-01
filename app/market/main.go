package main

import (
	"fmt"
	"log"
	"micros/app/market/subscriber"
	"micros/app/market/tasks"
	"micros/auth"
	"micros/config"
	"micros/database/redis"
	"micros/wrapper"

	_ "micros/broker/natsjs"

	_ "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/spf13/viper"

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

	a := auth.NewMicroAuth()

	service.Init(
		micro.Server(sgrpc.NewServer(server.Name(appName))),
		micro.Address(fmt.Sprintf(":%s", appPort)),
		micro.Auth(a),
		micro.WrapHandler(wrapper.NewRequestWrapper()),
		micro.WrapHandler(wrapper.NewAuthWrapper(a)))

	subscriber.Register(service)

	go tasks.ExecuteAll(service)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
