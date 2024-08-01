package orderbook

import (
	"fmt"
	"sort"
)

func NewBook(side OrderSide) *Book {
	var cmp = minCmp
	if side == BuySide {
		cmp = maxCmp
	}

	return &Book{
		levels: make(BookLevels, 128),
		cmp:    cmp,
	}
}

type BookLevels []*PriceLevel

type Book struct {
	cmp    func(a, b Price) bool
	levels BookLevels
}

func (b BookLevels) Make(order *Order) error {
	for _, pl := range b {
		if pl.price == order.Price {
			pl.orders = append(pl.orders, order)
			return nil
		}
	}

	b = append(b, NewPriceLevel(order.Price))

	// TODO: use cmp
	sort.Slice(b, func(i, j int) bool {
		return (b)[i].price < (b)[j].price
	})

	return nil
}

func (b BookLevels) Take(size Size) error {
	if len(b) == 0 {
		return fmt.Errorf("not enough liquidity in book")
	}

	var qtyLeft = size
	for _, pl := range b {
		if qtyLeft == 0 {
			return nil
		}

		if qtyLeft >= pl.totalSize {
			b = (b)[1:]
			qtyLeft -= pl.totalSize
			continue
		}

		for _, order := range pl.orders {
			qtyLeft -= order.Size
			pl.orders = pl.orders[1:]
		}
	}

	return nil
}

func (b BookLevels) Depth() int {
	return len(b)
}

func (b BookLevels) Top() (Price, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("no orders in book")
	}

	return b[0].price, nil
}

func minCmp(a, b Price) bool {
	return a < b
}

func maxCmp(a, b Price) bool {
	return a > b
}
