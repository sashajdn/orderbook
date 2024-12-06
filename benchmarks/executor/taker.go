package executor

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/sashajdn/orderbook/benchmarks/client"
	"github.com/sashajdn/orderbook/lob"
)

type TakerConfig struct {
	Users  uint
	Client client.Client
}

func NewTaker(config TakerConfig) *Taker {
	return &Taker{
		client: config.Client,
		users:  make(map[uint]struct{}, config.Users),
	}
}

var _ Executor = &Taker{}

type Taker struct {
	client client.Client
	users  map[uint]struct{}
}

func (t *Taker) RunIteration(ctx context.Context) error {
	for user := range t.users {
		if err := t.runIteration(ctx); err != nil {
			return fmt.Errorf("taker run iteration for user %d: %w", user, err)
		}

		return nil
	}

	return fmt.Errorf("invalid maker; no users set")
}

func (t *Taker) Name() string { return "taker" }

func (t *Taker) runIteration(ctx context.Context) error {
	side := lob.BuySide // TODO:
	if rand.Float64() < 0.5 {
		side = lob.SellSide
	}

	order := client.AddOrderRequest{
		OrderType: lob.MarketOrder,
		OrderSide: side,
		Size:      1,
	}

	if _, err := t.client.AddOrder(ctx, order); err != nil {
		return fmt.Errorf("add order: %w", err)
	}

	return nil
}
