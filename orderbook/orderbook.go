package orderbook

import (
	"fmt"
	"sync/atomic"
)

func NewOrderbook(size uint64) *Orderbook {
	return &Orderbook{
		asks: make([]*PriceLevel, 0, size),
		bids: make([]*PriceLevel, 0, size),
	}
}

type Orderbook struct {
	asks    BookLevels
	bids    BookLevels
	orderID uint64
}

func (o *Orderbook) Mid() (float64, error) {
	if len(o.asks) == 0 || len(o.bids) == 0 {
		return 0, fmt.Errorf(`no orders in book`)
	}

	return (o.asks[0].price + o.bids[0].price) / 2, nil
}

func (o *Orderbook) BestAsk() (float64, error) {
	if len(o.asks) == 0 {
		return 0, fmt.Errorf("no asks in book")
	}

	return o.asks[0].price, nil
}

func (o *Orderbook) BestBid() (float64, error) {
	if len(o.bids) == 0 {
		return 0, fmt.Errorf("no bids in book")
	}

	return o.bids[0].price, nil
}

func (o *Orderbook) Depth() (int, error) {
	return max(len(o.asks), len(o.bids)), nil
}

func (o *Orderbook) PlaceOrder(order *Order) (uint64, error) {
	if err := order.Validate(); err != nil {
		return 0, fmt.Errorf("invalid order: %w", err)
	}

	orderID := atomic.AddUint64(&o.orderID, 1)

	if order.OrderType == MarketOrder {
		switch order.Side {
		case BuySide:
			if err := o.asks.Take(order.Size); err != nil {
				return 0, fmt.Errorf("take order from asks: %w", err)
			}

			return orderID, nil
		case SellSide:
			if err := o.asks.Take(order.Size); err != nil {
				return 0, fmt.Errorf(`take order from bids: %w`, err)
			}

			return orderID, nil
		}
	}

	switch order.Side {
	case BuySide:
		if err := o.bids.Make(order); err != nil {
			return 0, fmt.Errorf("make order on bids: %w", err)
		}

		return orderID, nil
	case SellSide:
		if err := o.asks.Make(order); err != nil {
			return 0, fmt.Errorf(`make order on asks: %w`, err)
		}

		return orderID, nil
	}

	return 0, fmt.Errorf("invalid order")
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
