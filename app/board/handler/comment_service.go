package handler

import (
	"context"
	"fmt"

	"micros/app/board/models"
	"micros/database/mysql"
	boardV1 "micros/proto/board/v1"

	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
)

type CommentService struct{}

func (s *CommentService) GetAll(
	ctx context.Context,
	req *boardV1.CommentServiceGetAllRequest,
	rsp *boardV1.CommentServiceGetAllResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	db := mysql.Get()
	var data []*boardV1.Comment
	if result := db.Table("comments").
		Where("post_id = ?", req.PostId).
		Order("id desc").
		Find(&data); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	rsp.Data = data
	return nil
}

func (s *CommentService) Create(
	ctx context.Context,
	req *boardV1.CommentServiceCreateRequest,
	rsp *boardV1.CommentServiceCreateResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")
	db := mysql.Get()

	comment := &models.Comment{
		UserId:  userId,
		PostId:  req.PostId,
		Content: req.Content,
	}

	if result := db.Create(&comment); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	rsp.Data = &boardV1.Comment{
		Id:      fmt.Sprint(comment.ID),
		UserId:  comment.UserId,
		PostId:  comment.PostId,
		Content: comment.Content,
	}

	return nil
}
