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
	usersMap := make(map[uint64]struct{}, config.Users)
	for userID := 1; userID <= int(config.Users); userID++ {
		usersMap[uint64(userID)] = struct{}{}
	}

	return &Taker{
		client: config.Client,
		users:  usersMap,
	}
}

var _ Executor = &Taker{}

type Taker struct {
	client client.Client
	users  map[uint64]struct{}
}

func (t *Taker) RunIteration(ctx context.Context) error {
	for user := range t.users {
		if err := t.runIteration(ctx); err != nil {
			return fmt.Errorf("taker run iteration for user %d: %w", user, err)
		}

		return nil
	}

	return fmt.Errorf("invalid taker; no users set")
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
