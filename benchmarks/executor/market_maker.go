package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/sashajdn/orderbook/benchmarks/client"
)

var _ Executor = &MarketMaker{}

type MarketMakerOrderState struct {
	Pair        string
	BuyOrderID  string
	SellOrderID string
	mu          sync.Mutex
}

func NewMarketMaker(users int, client client.Client) *MarketMaker {
	return &MarketMaker{
		ordersMap: make(map[uint64]*MarketMakerOrderState, users),
		client:    client,
	}
}

type MarketMaker struct {
	ordersMap map[uint64]*MarketMakerOrderState
	client    client.Client
}

func (m *MarketMaker) RunIteration(ctx context.Context) error {
	for _, orders := range m.ordersMap {
		orders.mu.Lock()
		defer orders.mu.Unlock()

		if err := m.runIteration(ctx, orders); err != nil {
			return fmt.Errorf("run iteration: %w", err)
		}

		return nil
	}

	return fmt.Errorf("no orders set")
}

func (m *MarketMaker) runIteration(ctx context.Context, orderState *MarketMakerOrderState) error {
	if orderState.BuyOrderID == "" && orderState.SellOrderID == "" {
		var err error
		if orderState.BuyOrderID, err = m.addOrder(ctx, true); err != nil {
			return fmt.Errorf("buy side add order: %w", err)
		}

		if orderState.SellOrderID, err = m.addOrder(ctx, false); err != nil {
			return fmt.Errorf(`sell side add order: %w`, err)
		}

		return nil
	}

	// Randomise on editOrder ratio

	return nil
}

func (m *MarketMaker) addOrder(_ context.Context, _ bool) (string, error) {
	return "", fmt.Errorf("unimplemented")
}

func (m *MarketMaker) cancelOrder(_ context.Context) error {
	return fmt.Errorf("unimplemented")
}

func (m *MarketMaker) editOrder(_ context.Context) error {
	return fmt.Errorf("unimplemented")
}

func (m *MarketMaker) generateOrder(_ context.Context) error {
	return fmt.Errorf("unimplemented")
}
