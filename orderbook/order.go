package orderbook

import "fmt"

type OrderType byte

const (
	LimitOrder OrderType = iota + 1
	MarketOrder
)

type OrderSide byte

const (
	BuySide OrderSide = iota + 1
	SellSide
)

type Order struct {
	OrderType OrderType
	Side      OrderSide
	Price     Price
	Size      Size
}

func (o *Order) Validate() error {
	if o == nil {
		return fmt.Errorf("invalid order; nil")
	}

	if o.Size == 0 {
		return fmt.Errorf(`invalid order; zero size`)
	}

	return nil
}
