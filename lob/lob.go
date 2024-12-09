package lob

import (
	"fmt"
	"log/slog"
)

func NewOrderbook(size uint64) *Orderbook {
	return &Orderbook{
		asks:      NewBook(SellSide),
		bids:      NewBook(BuySide),
		sequencer: NewSequencer(),
	}
}

type Orderbook struct {
	asks      *Book
	bids      *Book
	orderID   uint64
	sequencer *Sequencer
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

func (o *Orderbook) Volume() (Size, Size) {
	return o.bids.levels.TotalVolume(), o.asks.levels.TotalVolume()
}

func (o *Orderbook) PlaceOrder(order *Order) (uint64, error) {
	if err := order.Validate(); err != nil {
		return 0, fmt.Errorf("invalid order: %w", err)
	}

	sequencedOrder := o.sequencer.Stamp(order)
	slog.Debug("LOB: placing order", "order", sequencedOrder.String())

	if order.OrderType == MarketOrder {
		switch order.Side {
		case BuySide:
			if _, err := o.asks.Take(sequencedOrder.Size); err != nil {
				return 0, fmt.Errorf("take order from asks: %w", err)
			}

			return sequencedOrder.ID, nil
		case SellSide:
			if _, err := o.bids.Take(sequencedOrder.Size); err != nil {
				return 0, fmt.Errorf(`take order from bids: %w`, err)
			}

			return sequencedOrder.ID, nil
		}
	}

	switch order.Side {
	case BuySide:
		o.bids.Make(sequencedOrder)
		return sequencedOrder.ID, nil
	case SellSide:
		o.asks.Make(sequencedOrder)
		return sequencedOrder.ID, nil
	}

	return 0, fmt.Errorf("invalid order")
}

func (o *Orderbook) CancelOrder(orderID uint64) error {
	return fmt.Errorf("unimplemented")
}

func (o *Orderbook) EditOrder(order *Order) error {
	return fmt.Errorf("unimplemented")
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
