package subscriber

import (
	"context"
	"fmt"
	"micros/app/order/models"
	"micros/queue"
)

func validate(ctx context.Context, microId string) error {
	isExist, err := new(models.DepositOrderEvent).Exist(microId)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}

	if isExist {
		return queue.ErrMessageConflicted
	}

	return nil
}
