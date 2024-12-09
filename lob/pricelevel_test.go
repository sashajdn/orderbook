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

	tests := []struct {
		name                   string
		ordersToAppend         []*Order
		ordersToTake           []*Order
		expectedRemainingSizes []Size
	}{
		{
			name: `sanity_check`,
			ordersToAppend: []*Order{
				{
					Size:          2,
					remainingSize: 2,
				},
			},
			ordersToTake: []*Order{
				{
					Size:          1,
					remainingSize: 1,
				},
			},
			expectedRemainingSizes: []Size{0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.True(t, len(tt.ordersToTake) == len(tt.expectedRemainingSizes), `test setup invalid; number of orders to take must be equal to expected remaining sizes`)

			pl := NewPriceLevel(1000.0)
			for _, order := range tt.ordersToAppend {
				pl.Append(order)
			}

			for i, order := range tt.ordersToTake {
				size, _ := pl.Take(order.Size)
				assert.Equal(t, tt.expectedRemainingSizes[i], size)
			}
		})
	}
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
