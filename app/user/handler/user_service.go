package handler

import (
	"context"
	"errors"
	"fmt"

	"micros/app/user/event"
	"micros/app/user/models"
	"micros/auth"
	"micros/auth/hash"
	userV1 "micros/proto/user/v1"

	"go-micro.dev/v4"
	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func NewUserService(service micro.Service, db *gorm.DB) *UserService {
	return &UserService{service: service, db: db}
}

type UserService struct {
	service micro.Service
	db      *gorm.DB
}

func (g *UserService) Register(
	ctx context.Context,
	req *userV1.RegisterRequest,
	rsp *userV1.RegisterResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	user := &models.User{Email: req.Email}

	if result := g.db.Where(user).First(user); result.RowsAffected > 0 {
		return microErrors.BadRequest("123", "user existed")
	}

	hashed, _ := hash.Make(req.Password)
	user.Password = hashed
	user.Name = req.Name

	if result := g.db.Create(user); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	event.UserCreated{Client: g.service.Client()}.Dispatch(fmt.Sprint(user.ID))

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

	var user models.User
	if result := g.db.Where(&models.User{Email: req.Email}).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
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

	var user models.User
	if result := g.db.First(&user, userId); errors.Is(result.Error, gorm.ErrRecordNotFound) {
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

func (g *UserService) GetList(
	ctx context.Context,
	req *userV1.GetListRequest,
	rsp *userV1.GetListResponse,
) (err error) {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var users []*userV1.User
	if result := g.db.Table("users").Where("id IN ?", req.UserId).Scan(&users); result.Error != nil {
		return microErrors.NotFound("123", result.Error.Error())
	}

	rsp.Result = &userV1.Result{Code: 200, Message: "success"}

	rsp.Data = users

	return nil
}
