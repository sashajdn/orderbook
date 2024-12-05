package orderbook

import (
	"fmt"
	"log/slog"
	"sync"
)

type FillStatus int8

const (
	Unfilled FillStatus = iota + 1
	Filled
	PartiallyFilled
)

func (f FillStatus) String() string {
	switch f {
	case Unfilled:
		return "unfilled"
	case Filled:
		return "filled"
	case PartiallyFilled:
		return "partially_filled"
	default:
		return "unknown"
	}
}

type FillEvent struct {
	Status  FillStatus
	Size    Size
	OrderID uint64
}

func (f FillEvent) String() string {
	return fmt.Sprintf("%d %s %.6f", f.OrderID, f.Status, f.Size)
}

type (
	Size  float64
	Price float64
)

func NewPriceLevel(price Price) *PriceLevel {
	return &PriceLevel{
		orderQueue: make([]*Order, 0),
	}
}

type PriceLevel struct {
	orderQueue []*Order
	price      Price
	totalSize  Size
	mu         sync.RWMutex
}

func (p *PriceLevel) Append(order *Order) {
	p.mu.Lock()
	defer p.mu.Unlock()

	slog.Debug("PL append: ", "pl price", fmt.Sprintf(`%.f6`, p.price), `order`, order.String())

	p.orderQueue = append(p.orderQueue, order)
	p.totalSize += order.Size
}

func (p *PriceLevel) Take(size Size) (Size, []*FillEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var (
		remainingSize = size
		fills         = make([]*FillEvent, 0, 1)
	)

	defer slog.Debug(`PL: matched`, "fills", fills)

	for _, order := range p.orderQueue {
		if remainingSize == 0 {
			break
		}

		switch {
		case remainingSize == order.remainingSize:
			fills = append(fills, &FillEvent{
				Status:  Filled,
				OrderID: order.ID,
				Size:    remainingSize,
			})

			order.remainingSize = 0
			remainingSize = 0
			p.orderQueue = p.orderQueue[1:]
			p.totalSize -= order.Size

			return remainingSize, fills
		case remainingSize < order.remainingSize:
			fills = append(fills, &FillEvent{
				Status:  PartiallyFilled,
				OrderID: order.ID,
				Size:    remainingSize,
			})

			order.remainingSize = order.remainingSize - remainingSize
			remainingSize = 0
			p.totalSize -= remainingSize

			return remainingSize, fills
		case remainingSize > order.remainingSize:
			fills = append(fills, &FillEvent{
				Status:  Filled,
				OrderID: order.ID,
				Size:    order.remainingSize,
			})

			remainingSize = remainingSize - order.remainingSize
			order.remainingSize = 0
			p.orderQueue = p.orderQueue[1:]
			p.totalSize -= order.Size
		}
	}

	return remainingSize, fills
}
