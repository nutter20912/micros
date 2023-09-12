package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	userV1 "micros/proto/user/v1"
	"os"
	"time"

	"github.com/go-micro/plugins/v4/auth/jwt/token"
	"github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/consul"

	"go-micro.dev/v4"
	"go-micro.dev/v4/auth"
	"go-micro.dev/v4/registry"
)

func main() {
	start := time.Now()
	call()
	log.Printf("Binomial took %s", time.Since(start))
}

func parse() {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiIiwic2NvcGVzIjpudWxsLCJtZXRhZGF0YSI6bnVsbCwiZXhwIjoxNjkzODM2NTU4LCJzdWIiOiJ0ZXN0In0.roAXCLDNWwmvJYBvGhemppaJAHPYsKEVFC5NSFu7m7khwhDiW_loO94pyBqBolYAEaEg7n-64K6jpb4a6GC7tmGjerxjlcDnZ8gwuuRIbkbSbJszE8pJj1i3OYyDP4EXxurSAECmmi0exfbYPaAL6Tg2_dppP9e4nnd4nhV2Nm1VJ4HQcIrgaAyRRi2v1Z4fRbhEbechPgQF7O0osgJWCXtsK_ByVA4vkHGs3RPiqmR55XKcwsEXmY-PinjJdB0_f52AWPrihcgPkY_KPbweeikntXIr2zh0ui4jNWyimeiFRsFOVqM4V3DbIcPVrD6aBsyewyFnfE01czPsYqpsOA"

	privKey, _ := os.ReadFile("private_key.pem")
	pubKey, _ := os.ReadFile("public_key.pem")

	j := token.New(
		token.WithPrivateKey(string(base64.StdEncoding.EncodeToString(privKey))),
		token.WithPublicKey(string(base64.StdEncoding.EncodeToString(pubKey))),
	)

	tok2, err := j.Inspect(tokenString)
	if err != nil {
		log.Fatalf("Inspect returned %v error, expected nil", err)
	}

	fmt.Println(tok2)
}

func gen() {
	privKey, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	j := token.New(
		token.WithPrivateKey(string(base64.StdEncoding.EncodeToString(privKey))),
	)

	t, err := j.Generate(&auth.Account{ID: "test"})
	if err != nil {
		log.Fatalf("Generate returned %v error, expected nil", err)
	}

	fmt.Println(t.Token)
}

func call() {
	re := consul.NewRegistry(registry.Addrs(":8500"))

	service := micro.NewService(
		micro.Name("User.Client"),
		micro.Registry(re),
	)
	service.Init()

	userService := userV1.NewUserService("srv.user", grpc.NewClient())
	rsp, err := userService.Login(context.Background(), &userV1.LoginRequest{
		Email:    "apple@gmail.com",
		Password: "zxczxc",
	})

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(rsp)
}
