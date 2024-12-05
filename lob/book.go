package orderbook

import (
	"fmt"
	"log/slog"
	"sort"
)

func NewBook(side OrderSide) *Book {
	var cmp = minCmp
	if side == BuySide {
		cmp = maxCmp
	}

	return &Book{
		levels: make(BookLevels, 0, 128),
		cmp:    cmp,
	}
}

type BookLevels []*PriceLevel

type Book struct {
	cmp    func(a, b Price) bool
	levels BookLevels
}

func (b *Book) Make(order *Order) {
	slog.Debug("Book make: ", "order", order.String())

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

	slog.Debug("PL not found, creating new PL @ ", "price", fmt.Sprintf("%.6f", order.Price))

	// If the price level is *not* found, then we append to end of the list & sort.
	b.levels = append(b.levels, NewPriceLevel(order.Price))

	sort.Slice(b.levels, func(i, j int) bool {
		return b.cmp(b.levels[i].price, b.levels[j].price)
	})
}

func (b *Book) Take(size Size) ([]*FillEvent, error) {
	priceLevels := b.levels
	if b.Depth() == 0 {
		return nil, fmt.Errorf("not enough liquidity in book")
	}

	var (
		qtyLeft    = size
		totalFills = []*FillEvent{}
	)
	for _, priceLevel := range priceLevels {
		if qtyLeft == 0 {
			break
		}

		var fills []*FillEvent
		qtyLeft, fills = priceLevel.Take(qtyLeft)
		totalFills = append(totalFills, fills...)
	}

	// TODO: we need a way to manage what happens when there isn't enough liqudity in the book
	// BUG: there is also the case here whereby, we want to be sure we have enough liqudity in the book before accepting an order in the
	//      the book

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
