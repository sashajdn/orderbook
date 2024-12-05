package lob

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPricelevel_SanityCheck(t *testing.T) {
	t.Parallel()

	pl := NewPriceLevel(1000.0)

	orders := generateOrders(10, 0, 1000.0, 0, []uint64{1, 2, 3})
	for _, order := range orders {
		fmt.Println("Adding order to price level: ", order)
		pl.Append(order)
		fmt.Println("PL total size: ", pl.totalSize)
	}

	assert.Len(t, pl.orderQueue, 10)

	var totalSize Size
	for _, order := range pl.orderQueue {
		require.True(t, order.Price == Price(1000.0))
		totalSize += order.Size
	}

	assert.True(t, totalSize == pl.totalSize)
}

func TestPricelevel_Take(t *testing.T) {
	t.Parallel()

	pl := NewPriceLevel(1000.0)

	orders := generateOrders(10, 0, 1000.0, 0, []uint64{1})
	for _, order := range orders {
		fmt.Println("Adding order to price level: ", order)
		pl.Append(order)
		fmt.Println("PL total size: ", pl.totalSize)
	}

	size, fills := pl.Take(2)
	assert.Len(t, fills, 2)
	assert.Equal(t, Size(0), size)
	assert.Equal(t, Size(8), pl.totalSize)
}

func generateOrders(n, m uint, midpoint, spread float64, sizeRange []uint64) []*Order {
	orders := make([]*Order, 0, n+m)

	for i := 0; i < int(n); i++ {
		price := (rand.NormFloat64() * spread / 2) + midpoint

		i := rand.Intn(len(sizeRange))
		size := sizeRange[i]

		order := defaultSequencer.NewOrder(LimitOrder, BuySide, Price(price), Size(size))
		orders = append(orders, order)
	}

	for j := 0; j < int(m); j++ {
		price := (rand.NormFloat64() * spread / 2) + midpoint

		i := rand.Intn(len(sizeRange))
		size := sizeRange[i]

		order := defaultSequencer.NewOrder(LimitOrder, SellSide, Price(price), Size(size))
		orders = append(orders, order)
	}

	return orders
}
