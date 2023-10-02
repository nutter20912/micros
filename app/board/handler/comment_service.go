package handler

import (
	"context"
	"encoding/base64"
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
	tx := db.Table("comments").Where("post_id = ?", req.PostId)

	if req.Cursor != nil && *req.Cursor != "" {
		cursor, err := base64.StdEncoding.DecodeString(*req.Cursor)
		if err != nil {
			return microErrors.InternalServerError("123", err.Error())
		}

		if string(cursor) != "" {
			tx.Where("id < ?", string(cursor))
			fmt.Println(string(cursor))
		}
	}

	var data []*boardV1.Comment
	if result := tx.Order("id desc").Limit(PAGE_LIMIT).Find(&data); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	next := ""
	if len(data) == PAGE_LIMIT {
		next = base64.StdEncoding.EncodeToString([]byte(data[len(data)-1].Id))
	}

	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	rsp.Data = data
	rsp.Paginator = &boardV1.Paginator{
		NextCursor: next,
	}

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
