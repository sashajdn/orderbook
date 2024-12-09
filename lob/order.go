package lob

import "fmt"

type OrderType byte

const (
	LimitOrder OrderType = iota + 1
	MarketOrder
)

func (o OrderType) String() string {
	switch o {
	case LimitOrder:
		return "limit_order"
	case MarketOrder:
		return "market_order"
	default:
		return "unknown"
	}
}

type OrderSide byte

const (
	BuySide OrderSide = iota + 1
	SellSide
)

func (o OrderSide) String() string {
	switch o {
	case BuySide:
		return "buy"
	case SellSide:
		return "sell"
	default:
		return "unknown"
	}
}

func NewOrder(orderType OrderType, side OrderSide, price Price, size Size) *Order {
	return &Order{
		OrderType:     orderType,
		Side:          side,
		Price:         price,
		Size:          size,
		remainingSize: size,
	}
}

type Order struct {
	OrderType     OrderType
	Side          OrderSide
	Price         Price
	Size          Size
	ID            uint64
	remainingSize Size
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

func (o *Order) String() string {
	return fmt.Sprintf(`%s @ %.6f : id=%d type=%s size=%.6f remsize=%.6f`, o.Side, o.Price, o.ID, o.OrderType, o.Size, o.remainingSize)
}
