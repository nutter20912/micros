package handler

import (
	"context"

	"micros/app/board/models"
	"micros/database/mysql"
	boardV1 "micros/proto/board/v1"

	microErrors "go-micro.dev/v4/errors"
)

type BoardService struct{}

func (s *BoardService) GetAll(
	ctx context.Context,
	req *boardV1.PostServiceGetAllRequest,
	rsp *boardV1.PostServiceGetAllResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	db := mysql.Get()
	var data []*boardV1.Post
	if result := db.Table("posts").Find(&data); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	rsp.Data = data

	return nil
}

func (s *BoardService) Create(
	ctx context.Context,
	req *boardV1.PostServiceCreateRequest,
	rsp *boardV1.PostServiceCreateResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	db := mysql.Get()
	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserId:  1,
	}

	if result := db.Create(&post); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	return nil
}

func (s *BoardService) Update(
	ctx context.Context,
	req *boardV1.PostServiceUpdateRequest,
	rsp *boardV1.PostServiceUpdateResponse,
) error {
	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	return nil
}

func (s *BoardService) Get(
	ctx context.Context,
	req *boardV1.PostServiceGetRequest,
	rsp *boardV1.PostServiceGetResponse,
) error {
	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	return nil
}
