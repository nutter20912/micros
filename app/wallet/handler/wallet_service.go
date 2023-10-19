package handler

import (
	"context"
	"encoding/json"

	"micros/app/wallet/models"
	walletV1 "micros/proto/wallet/v1"

	"go-micro.dev/v4"
	microErrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WalletService struct {
	Service micro.Service
}

func (s *WalletService) GetEvent(
	ctx context.Context,
	req *walletV1.GetEventRequest,
	rsp *walletV1.GetEventResponse,
) error {
	if err := req.Validate(); err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	userId, _ := metadata.Get(ctx, "user_id")

	events, paginatior, err := new(models.WalletEvent).Get(userId, req.Page, req.Limit)
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data []*walletV1.WalletEvent
	bytes, _ := json.Marshal(events)
	json.Unmarshal(bytes, &data)

	var p *walletV1.Paginator
	bytes, _ = json.Marshal(paginatior)
	json.Unmarshal(bytes, &p)

	rsp.Data = data
	rsp.Paginator = p

	return nil
}

func (s *WalletService) Get(
	ctx context.Context,
	req *emptypb.Empty,
	rsp *walletV1.GetResponse,
) error {
	userId, _ := metadata.Get(ctx, "user_id")

	wallet, err := new(models.Wallet).Get(userId)
	if err != nil {
		return microErrors.BadRequest("222", err.Error())
	}

	var data *walletV1.Wallet
	bytes, _ := json.Marshal(wallet)
	json.Unmarshal(bytes, &data)

	rsp.Data = data

	return nil
}
