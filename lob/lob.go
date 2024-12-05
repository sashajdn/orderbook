package orderbook

import (
	"fmt"
	"log/slog"
	"sync/atomic"
)

func NewOrderbook(size uint64) *Orderbook {
	return &Orderbook{
		asks: NewBook(SellSide),
		bids: NewBook(BuySide),
	}
}

type Orderbook struct {
	asks    *Book
	bids    *Book
	orderID uint64
}

func (o *Orderbook) Mid() (Price, error) {
	bbp, err := o.bids.Top()
	if err != nil {
		return 0, fmt.Errorf("fetch bids top: %w", err)

	}

	bap, err := o.asks.Top()
	if err != nil {
		return 0, fmt.Errorf("fetch bids top: %w", err)

	}

	return (bbp + bap) / 2, nil
}

func (o *Orderbook) BestAsk() (Price, error) {
	return o.asks.Top()
}

func (o *Orderbook) BestBid() (Price, error) {
	return o.bids.Top()
}

func (o *Orderbook) Depth() int {
	return max(o.asks.Depth(), o.bids.Depth())
}

func (o *Orderbook) PlaceOrder(order *Order) (uint64, error) {
	// Validate order.
	if err := order.Validate(); err != nil {
		return 0, fmt.Errorf("invalid order: %w", err)
	}

	orderID := atomic.AddUint64(&o.orderID, 1)
	order.ID = orderID
	slog.Debug(`LOB: placing order`, "order", order.String())

	if order.OrderType == MarketOrder {
		switch order.Side {
		case BuySide:
			if _, err := o.asks.Take(order.Size); err != nil {
				return 0, fmt.Errorf("take order from asks: %w", err)
			}

			return orderID, nil
		case SellSide:
			if _, err := o.asks.Take(order.Size); err != nil {
				return 0, fmt.Errorf(`take order from bids: %w`, err)
			}

			return orderID, nil
		}
	}

	switch order.Side {
	case BuySide:
		o.bids.Make(order)
		return orderID, nil
	case SellSide:
		o.asks.Make(order)
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
