package handler

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"micros/app/board/models"
	"micros/database/mysql"
	boardV1 "micros/proto/board/v1"

	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"gorm.io/gorm"
)

type PostService struct{}

var (
	PAGE_LIMIT = 4
)

func (s *PostService) GetAll(
	ctx context.Context,
	req *boardV1.PostServiceGetAllRequest,
	rsp *boardV1.PostServiceGetAllResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	db := mysql.Get()

	tx := db.Table("posts")

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

	var data []*boardV1.Post
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

func (s *PostService) Create(
	ctx context.Context,
	req *boardV1.PostServiceCreateRequest,
	rsp *boardV1.PostServiceCreateResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")

	db := mysql.Get()
	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserId:  userId,
	}

	if result := db.Create(&post); result.Error != nil {
		return microErrors.InternalServerError("123", result.Error.Error())
	}

	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	rsp.Data = &boardV1.Post{
		Id:      fmt.Sprint(post.ID),
		UserId:  post.UserId,
		Title:   post.Title,
		Content: post.Content,
	}

	return nil
}

func (s *PostService) Update(
	ctx context.Context,
	req *boardV1.PostServiceUpdateRequest,
	rsp *boardV1.PostServiceUpdateResponse,
) error {
	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	return nil
}

func (s *PostService) Get(
	ctx context.Context,
	req *boardV1.PostServiceGetRequest,
	rsp *boardV1.PostServiceGetResponse,
) error {

	db := mysql.Get()

	var post models.Post
	if result := db.First(&post, req.Id); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return microErrors.NotFound("123", "post not found")
	}

	rsp.Result = &boardV1.Result{Code: 200, Message: "success"}
	rsp.Data = &boardV1.Post{
		Id:      fmt.Sprint(post.ID),
		UserId:  post.UserId,
		Title:   post.Title,
		Content: post.Content,
	}
	return nil
}
