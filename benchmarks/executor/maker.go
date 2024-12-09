package executor

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/sashajdn/orderbook/benchmarks/client"
	"github.com/sashajdn/orderbook/lob"
)

type MakerConfig struct {
	Users       uint
	LaplaceBeta float64
	Midprice    lob.Price
	Spread      float64
	Client      client.Client
}

func NewMaker(config MakerConfig) *Maker {
	usersMap := make(map[uint64]struct{}, config.Users)
	for userID := 1; userID <= int(config.Users); userID++ {
		usersMap[uint64(userID)] = struct{}{}
	}

	return &Maker{
		client:      config.Client,
		users:       usersMap,
		laplaceBeta: config.LaplaceBeta,
		midprice:    config.Midprice,
		spread:      config.Spread,
	}
}

var _ Executor = &Maker{}

type Maker struct {
	client      client.Client
	users       map[uint64]struct{}
	laplaceBeta float64
	midprice    lob.Price
	spread      float64
}

func (m *Maker) RunIteration(ctx context.Context) error {
	for user := range m.users {
		if err := m.runIteration(ctx); err != nil {
			return fmt.Errorf("maker run iteration for user %d: %w", user, err)
		}

		return nil
	}

	return fmt.Errorf("invalid maker; no users set")
}

func (m *Maker) Name() string { return "maker" }

func (m *Maker) runIteration(ctx context.Context) error {
	side := lob.BuySide // TODO:
	if rand.Float64() < 0.5 {
		side = lob.SellSide
	}

	var size lob.Size = 1 // TODO:

	order, err := m.generateOrder(lob.Price(m.midprice), lob.Price(m.spread), side, size, m.laplaceBeta)
	if err != nil {
		return fmt.Errorf("generate order: %w", err)
	}

	if _, err = m.client.AddOrder(ctx, order); err != nil {
		return fmt.Errorf("add order: %w", err)
	}

	return nil
}

func (m *Maker) generateOrder(midprice lob.Price, spread lob.Price, side lob.OrderSide, size lob.Size, laplaceBeta float64) (client.AddOrderRequest, error) {
	delta := spread * 2 // This should put the price point sufficiently away from the midpoint.

	var price lob.Price
	switch side {
	case lob.BuySide:
		price = midprice + delta + lob.Price(laplaceRandom(laplaceBeta))

	case lob.SellSide:
		price = midprice - delta - lob.Price(laplaceRandom(laplaceBeta))
	}

	return client.AddOrderRequest{
		OrderType: lob.LimitOrder,
		OrderSide: side,
		Price:     price,
		Size:      size,
	}, nil
}

func laplaceRandom(beta float64) float64 {
	u := rand.Float64() - 0.5
	if u < 0 {
		return math.Log(1-2*math.Abs(u)) * beta
	}

	return -math.Log(1-2*math.Abs(u)) * beta
}
