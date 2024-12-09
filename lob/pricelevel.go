package lob

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
	Price   Price
	Size    Size
	OrderID uint64
}

func (f FillEvent) String() string {
	return fmt.Sprintf(`%d %s %.6f %.6f`, f.OrderID, f.Status, f.Price, f.Size)
}

type (
	Size  float64
	Price float64
)

func NewPriceLevel(price Price) *PriceLevel {
	return &PriceLevel{
		orderQueue: make([]*Order, 0),
		price:      price,
	}
}

type PriceLevel struct {
	orderQueue []*Order
	price      Price
	totalSize  Size
	mu         sync.RWMutex
}

func (p *PriceLevel) String() string {
	return fmt.Sprintf("PL: price=%.6f size=%.6f", p.price, p.totalSize)
}

func (p *PriceLevel) Append(order *Order) {
	p.mu.Lock()
	defer p.mu.Unlock()

	slog.Debug("Appending order to pricelevel: ", "order", order.String(), "pricelevel", p.String())

	p.orderQueue = append(p.orderQueue, order)
	p.totalSize += order.Size

}

func (p *PriceLevel) Take(size Size) (Size, []*FillEvent) {
	if size == 0 {
		return 0, []*FillEvent{}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	var (
		remainingSize = size
		fills         = make([]*FillEvent, 0, 1)
	)

	defer slog.Debug("PL: matched taker orders", "fills", fills)

	for _, order := range p.orderQueue {
		if remainingSize == 0 {
			break
		}

		switch {
		case remainingSize == order.remainingSize:
			fills = append(fills, &FillEvent{
				Status:  Filled,
				OrderID: order.ID,
				Price:   p.price,
				Size:    remainingSize,
			})

			p.totalSize -= remainingSize
			order.remainingSize = 0

			p.orderQueue = p.orderQueue[1:]

			return 0, fills
		case remainingSize < order.remainingSize:
			fills = append(fills, &FillEvent{
				Status:  PartiallyFilled,
				Price:   p.price,
				OrderID: order.ID,
				Size:    remainingSize,
			})

			order.remainingSize -= -remainingSize
			p.totalSize -= remainingSize

			return 0, fills
		case remainingSize > order.remainingSize:
			fills = append(fills, &FillEvent{
				Status:  Filled,
				Price:   p.price,
				OrderID: order.ID,
				Size:    order.remainingSize,
			})

			p.totalSize -= order.remainingSize
			remainingSize -= order.remainingSize
			order.remainingSize = 0

			p.orderQueue = p.orderQueue[1:]
		}
	}

	return remainingSize, fills
}

func (p *PriceLevel) Volume() Size {
	return p.totalSize
}

func (p *PriceLevel) NumberOfOrders() int {
	return len(p.orderQueue)
}
