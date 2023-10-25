package subscriber

import (
	"context"
	"fmt"
	"micros/app/wallet/models"
	baseEvent "micros/event"
)

func validate(ctx context.Context, microId string) error {
	isExist, err := new(models.WalletEvent).Exist(microId)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}

	if isExist {
		return baseEvent.ErrMessageConflicted
	}

	return nil
}
