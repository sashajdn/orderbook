package lob

import (
	"sync"
	"sync/atomic"
)

var defaultSequencer = &Sequencer{}

func NewSequencer() *Sequencer {
	return &Sequencer{}
}

type Sequencer struct {
	mu sync.Mutex
	id uint64
}

func (s *Sequencer) NewOrder(orderType OrderType, side OrderSide, price Price, size Size) *Order {
	orderID := s.generateNextID()
	return &Order{
		OrderType:     orderType,
		Side:          side,
		Price:         price,
		Size:          size,
		ID:            orderID,
		remainingSize: size,
	}
}

func (s *Sequencer) Stamp(order *Order) *Order {
	order.ID = s.generateNextID()
	return order
}

func (s *Sequencer) generateNextID() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return atomic.AddUint64(&s.id, 1)
}
