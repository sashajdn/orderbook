package lob

import (
	"sync"
	"sync/atomic"
)

var defaultSequencer = &Sequencer{}

type Sequencer struct {
	mu sync.Mutex
	id uint64
}

func (s *Sequencer) NewOrder(orderType OrderType, side OrderSide, price Price, size Size) *Order {
	orderID := atomic.AddUint64(&s.id, 1)
	return &Order{
		OrderType:     orderType,
		Side:          side,
		Price:         price,
		Size:          size,
		ID:            orderID,
		remainingSize: size,
	}
}
