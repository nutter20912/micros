package queue

import (
	"context"
	"fmt"
	"micros/event"
	"strconv"
	"time"

	"go-micro.dev/v4/metadata"
)

func getMedata(ctx context.Context) (metadata.Metadata, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, ErrNoMetadataReceived
	}
	return md, nil
}

func MicroId(ctx context.Context) (string, error) {
	md, err := getMedata(ctx)
	if err != nil {
		return "", err
	}

	val, ok := md["Micro-Id"]
	if !ok {
		return "", fmt.Errorf("Micro-Id not found")
	}

	return val, nil
}
func Validate(ctx context.Context) error {
	md, err := getMedata(ctx)
	if err != nil {
		return err
	}

	if ttl, ok := md[event.PUB_OPTIONS_TTL]; ok {
		val, err := strconv.ParseInt(ttl, 10, 64)
		if err != nil {
			return fmt.Errorf("wrong time")
		}

		if val < time.Now().Unix() {
			return ErrMessageExpired
		}
	}

	return nil
}

func ErrReportOrIgnore(err error) error {
	switch err {
	case ErrMessageExpired:
		fmt.Println("ErrMessageExpired")
		return nil
	case ErrMessageConflicted:
		fmt.Println("ErrMessageConflicted")
		return nil
	default:
		return err
	}
}

type EventModel interface {
	Exist(string) (bool, error)
}

func CheckMsgId(model EventModel, microId string) error {
	isExist, err := model.Exist(microId)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}

	if isExist {
		return ErrMessageConflicted
	}

	return nil
}
