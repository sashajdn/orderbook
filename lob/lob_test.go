package orderbook

import (
	"fmt"
	"testing"

	"github.com/sashajdn/orderbook/pkg/slog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLOB_SanityCheck(t *testing.T) {
	t.Parallel()
	slog.Init(slog.Debug)

	lob := NewOrderbook(128)

	// Add liquidity to the book.
	addSymmetricalDepthOf3(t, lob)

	idBuy, err := lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      BuySide,
		Size:      1,
	})
	require.NoError(t, err)

	idSell, err := lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      SellSide,
		Size:      1,
	})
	require.NoError(t, err)

	assert.True(t, idSell > idBuy)
	assert.Equal(t, 3, lob.Depth())

	_, err = lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      BuySide,
		Size:      2,
	})
	require.NoError(t, err)

	_, err = lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      SellSide,
		Size:      2,
	})
	require.NoError(t, err)

	mid, err := lob.Mid()
	require.NoError(t, err)

	assert.Equal(t, 1, lob.Depth())
	assert.Equal(t, 1000, mid)

	ba, err := lob.BestAsk()
	require.NoError(t, err)

	bb, err := lob.BestBid()
	require.NoError(t, err)

	spread := ba - bb
	assert.Equal(t, 6, spread)
}

func addSymmetricalDepthOf3(t *testing.T, lob *Orderbook) {
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
		Price:     1001,
		Size:      1,
	})
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
		Price:     1001,
		Size:      2,
	})
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
		Price:     1002,
		Size:      1,
	})
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
		Price:     1003,
		Size:      1,
	})
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     999,
		Size:      1,
	})
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     999,
		Size:      2,
	})
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     998,
		Size:      1,
	})
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     997,
		Size:      1,
	})

	// Print book
	printBook(t, lob)

	// Assertions
	assert.Equal(t, 3, lob.Depth())

	bb, err := lob.BestBid()
	require.NoError(t, err)

	ba, err := lob.BestAsk()
	require.NoError(t, err)

	mid, err := lob.Mid()
	require.NoError(t, err)

	assert.Equal(t, 999, bb)
	assert.Equal(t, 1001, ba)
	assert.Equal(t, 1000, mid)
}

func printBook(_ *testing.T, lob *Orderbook) {
	bids := lob.bids

	for i := len(lob.asks.levels) - 1; i >= 0; i-- {
		ask := lob.asks.levels[i]
		fmt.Println("Ask: ", i, ask.totalSize, ask.price)
	}

	fmt.Println()

	for i, bid := range bids.levels {
		fmt.Println("bid: ", i, bid.totalSize, bid.price)
	}
}
