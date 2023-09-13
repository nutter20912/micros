package handler

import (
	"context"
	"errors"
	"fmt"

	"micros/app/user/models"
	"micros/auth"
	"micros/auth/hash"
	"micros/database/mysql"
	userV1 "micros/proto/user/v1"

	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserService struct{}

func (g *UserService) Register(
	ctx context.Context,
	req *userV1.RegisterRequest,
	rsp *userV1.RegisterResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	db := mysql.Get()
	user := &models.User{Email: req.Email}

	if result := db.Where(user).First(user); result.RowsAffected > 0 {
		return microErrors.BadRequest("123", "user existed")
	}

	hashed, _ := hash.Make(req.Password)
	user.Password = hashed
	user.Name = req.Name

	if result := db.Create(user); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	rsp.Result = &userV1.Result{Code: 200, Message: "success"}

	return nil
}

func (g *UserService) Login(
	ctx context.Context,
	req *userV1.LoginRequest,
	rsp *userV1.LoginResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	db := mysql.Get()

	var user models.User
	if result := db.Where(&models.User{Email: req.Email}).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return microErrors.NotFound("123", "user not found")
	}

	if !hash.Check(req.Password, user.Password) {
		return microErrors.BadRequest("123", "auth failed")
	}

	token, _ := auth.Generate(fmt.Sprint(user.ID))

	rsp.Result = &userV1.Result{Code: 200, Message: "success"}
	rsp.Data = &userV1.LoginResponse_Data{Token: token.Token}

	return nil
}

func (g *UserService) Update(
	ctx context.Context,
	req *userV1.UpdateRequest,
	rsp *userV1.UpdateResponse,
) (err error) {

	rsp.Result = &userV1.Result{
		Code:    200,
		Message: "success",
	}
	return nil
}

func (g *UserService) Get(
	ctx context.Context,
	req *emptypb.Empty,
	rsp *userV1.GetResponse,
) (err error) {
	userId, _ := metadata.Get(ctx, "user_id")

	db := mysql.Get()

	var user models.User
	if result := db.First(&user, userId); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return microErrors.NotFound("123", "user not found")
	}

	rsp.Result = &userV1.Result{Code: 200, Message: "success"}
	rsp.Data = &userV1.User{
		Id:        userId,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.String(),
	}

	return nil
}
