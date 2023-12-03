package event

import (
	"context"
	"fmt"
	"time"

	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
)

var (
	PUB_OPTIONS_TTL = "TTL"
)

var (
	ORDER_DEPOSIT_CREATED = "order.deposit.created"
	ORDER_SPOT_CREATED    = "order.spot.created"
	ORDER_SPOT_UPDATED    = "order.spot.updated"

	MARKET_MATCHED       = "market.spot.matched"
	MARKET_PRICE_CHANGED = "market.price.changed"

	WALLET_TRANSACTION     = "wallet.transaction"
	WALLET_BALANCE_CHECKED = "wallet.balance.checked"

	USER_CREATED = "user.created"
)

func New(c client.Client) *Event {
	return &Event{client: c}
}

type Event struct {
	client client.Client
}

type Publisher interface {
	Publish(context.Context, client.Client) error
}

func (e *Event) Dispatch(
	ctx context.Context,
	p Publisher,
	opts ...DispatchOption,
) error {
	mdOpts := map[string]string{}

	for _, o := range opts {
		o(mdOpts)
	}

	ctx = metadata.NewContext(ctx, mdOpts)

	if err := p.Publish(ctx, e.client); err != nil {
		return fmt.Errorf("publish error: %v", err)
	}

	return nil
}

type DispatchOptions map[string]string

type DispatchOption func(DispatchOptions)

func SetTTL(ttl time.Duration) DispatchOption {
	return func(o DispatchOptions) {
		o[PUB_OPTIONS_TTL] = fmt.Sprint(time.Now().Add(ttl).Unix())
	}
}
