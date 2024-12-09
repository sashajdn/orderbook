package lob

import (
	"fmt"
	"log/slog"
	"sort"
	"sync"
)

func NewBook(side OrderSide) *Book {
	var cmp = minCmp
	if side == BuySide {
		cmp = maxCmp
	}

	return &Book{
		side:   side,
		levels: make(BookLevels, 0, 128),
		cmp:    cmp,
	}
}

type BookLevels []*PriceLevel

func (b BookLevels) TotalVolume() Size {
	var totalVolume Size
	for _, pl := range b {
		totalVolume += pl.totalSize
	}

	return totalVolume
}

type Book struct {
	side   OrderSide
	cmp    func(a, b Price) bool
	levels BookLevels
	mu     sync.RWMutex
}

func (b *Book) Side() OrderSide {
	return b.side
}

func (b *Book) TotalVolume() Size {
	return b.levels.TotalVolume()
}

func (b *Book) Make(order *Order) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Iterate through price levels until we find the price level we want; append order.
	for _, pl := range b.levels {
		if pl == nil {
			break
		}

		if pl.price == order.Price {
			pl.Append(order)
			return
		}
	}

	slog.Debug("Creating new pricelevel @", "price", fmt.Sprintf("%.6f", order.Price))

	// If the price level is *not* found, then we append to end of the list & sort.
	pl := NewPriceLevel(order.Price)
	pl.Append(order)
	b.levels = append(b.levels, pl)

	sort.Slice(b.levels, func(i, j int) bool {
		return b.cmp(b.levels[i].price, b.levels[j].price)
	})
}

func (b *Book) Take(size Size) ([]*FillEvent, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	priceLevels := b.levels
	if b.Depth() == 0 {
		// TODO: we should store this per the price level rather than being in a position whereby we need to calculate.
		return nil, fmt.Errorf("not enough liquidity in book %.2f/%.2f", size, b.levels.TotalVolume())
	}

	var (
		qtyLeft    = size
		totalFills = []*FillEvent{}
	)

	toRemoveFrom := -1
	for i, priceLevel := range priceLevels {
		if qtyLeft == 0 {
			break
		}

		var fills []*FillEvent
		qtyLeft, fills = priceLevel.Take(qtyLeft)
		totalFills = append(totalFills, fills...)

		if priceLevel.Volume() == 0 {
			toRemoveFrom = max(toRemoveFrom, i)
		}

		if qtyLeft <= 0 {
			break
		}
	}

	// Clean up price levels
	if toRemoveFrom >= 0 {
		if toRemoveFrom >= len(b.levels) {
			b.levels = make(BookLevels, 0, 128)
		} else {
			b.levels = b.levels[toRemoveFrom+1:]
		}
	}

	return totalFills, nil
}

func (b *Book) Depth() int {
	return len(b.levels)
}

func (b *Book) Top() (Price, error) {
	if len(b.levels) == 0 {
		return 0, fmt.Errorf("no orders in book")
	}

	return b.levels[0].price, nil
}

func minCmp(a, b Price) bool {
	return a < b
}

func maxCmp(a, b Price) bool {
	return a > b
}
