package wrapper

import (
	"context"
	"errors"
	"log"
	"micros/auth"

	baseAuth "go-micro.dev/v4/auth"
	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
)

func NewRequestWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			l := make(map[string]interface{})

			err := fn(ctx, req, rsp)

			l["request"] = req.Method()

			if err != nil {
				l["error"] = err.Error()
			}

			log.Printf("[request_log]:  %v", l)

			return err
		}
	}
}

func NewAuthWrapper(a baseAuth.Auth) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			switch req.Method() {
			case "UserService.Register":
			case "UserService.Login":
			default:
				md, ok := metadata.FromContext(ctx)
				if !ok {
					return errors.New("no metadata found")
				}

				acc, err := auth.VerifyToken(a, md)
				if err != nil {
					return microErrors.Unauthorized("401", err.Error())
				}

				if err := auth.VerifyAbility(req, acc); err != nil {
					return err
				}

				ctx = metadata.Set(ctx, "user_id", acc.ID)
			}

			return fn(ctx, req, rsp)
		}
	}
}
