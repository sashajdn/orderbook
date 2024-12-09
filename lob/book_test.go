package lob

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		side             OrderSide
		ordersToPlace    []*Order
		exceptedDepth    int
		expectedLiqudity Size
	}{
		{
			name: "sanity_check_maker_buy_side",
			side: BuySide,
			ordersToPlace: []*Order{
				{
					OrderType:     LimitOrder,
					Price:         1000,
					Size:          1,
					remainingSize: 1,
				},
				{
					OrderType:     LimitOrder,
					Price:         999,
					Size:          1,
					remainingSize: 1,
				},
				{
					OrderType:     LimitOrder,
					Price:         1001,
					Size:          1,
					remainingSize: 1,
				},
				{
					OrderType:     LimitOrder,
					Price:         1001,
					Size:          4,
					remainingSize: 4,
				},
			},
			exceptedDepth:    3,
			expectedLiqudity: 7,
		},
		{
			name: "maker_followed_by_taker",
			side: BuySide,
			ordersToPlace: []*Order{
				{
					OrderType:     LimitOrder,
					Price:         1000,
					Size:          1,
					remainingSize: 1,
				},
				{
					OrderType:     LimitOrder,
					Price:         999,
					Size:          10,
					remainingSize: 10,
				},
				{
					OrderType:     LimitOrder,
					Price:         998,
					Size:          2,
					remainingSize: 2,
				},
				{
					OrderType:     MarketOrder,
					Size:          1,
					remainingSize: 1,
				},
				{
					OrderType:     MarketOrder,
					Size:          11,
					remainingSize: 11,
				},
			},
			exceptedDepth:    1,
			expectedLiqudity: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			book := NewBook(BuySide)

			// Place orders.
			for _, order := range tt.ordersToPlace {
				switch order.OrderType {
				case MarketOrder:
					_, err := book.Take(order.Size)
					require.NoError(t, err)
				case LimitOrder:
					book.Make(order)
				}
			}

			assert.Equal(t, tt.exceptedDepth, book.Depth())
			assert.Equal(t, tt.expectedLiqudity, book.levels.TotalVolume())
		})
	}
}
