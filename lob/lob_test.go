package lob

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

	fmt.Println("========= Placing Market order Buy side of size 1")

	idBuy, err := lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      BuySide,
		Size:      1,
	})
	require.NoError(t, err)

	printBook(t, lob)

	fmt.Println(`========= Placing Market order Sell side of size 1`)

	idSell, err := lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      SellSide,
		Size:      1,
	})
	require.NoError(t, err)

	printBook(t, lob)

	assert.True(t, idSell > idBuy)
	assert.Equal(t, 3, lob.Depth())

	fmt.Println("========= Placing Market order Buy side to remove 1st depth")

	_, err = lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      BuySide,
		Size:      3,
	})
	require.NoError(t, err)

	fmt.Println(`========= Placing Market order Sell side to remove 1st depth`)

	_, err = lob.PlaceOrder(&Order{
		OrderType: MarketOrder,
		Side:      SellSide,
		Size:      3,
	})
	require.NoError(t, err)

	printBook(t, lob)

	fmt.Println("========")

	mid, err := lob.Mid()
	require.NoError(t, err)

	assert.Equal(t, 1, lob.Depth())
	assert.Equal(t, Price(1000), mid)

	ba, err := lob.BestAsk()
	require.NoError(t, err)

	bb, err := lob.BestBid()
	require.NoError(t, err)

	spread := ba - bb
	assert.Equal(t, 6, spread)
}

func addSymmetricalDepthOf3(t *testing.T, lob *Orderbook) {
	// Print book
	printBook(t, lob)

	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     1001,
		Size:      1,
	})

	// Print book
	printBook(t, lob)

	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     1001,
		Size:      2,
	})

	// Print book
	printBook(t, lob)
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     1002,
		Size:      1,
	})

	// Print book
	printBook(t, lob)
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      SellSide,
		Price:     1003,
		Size:      1,
	})

	// Print book
	printBook(t, lob)
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
		Price:     999,
		Size:      1,
	})

	// Print book
	printBook(t, lob)
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
		Price:     999,
		Size:      2,
	})

	// Print book
	printBook(t, lob)
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
		Price:     998,
		Size:      1,
	})

	// Print book
	printBook(t, lob)
	lob.PlaceOrder(&Order{
		OrderType: LimitOrder,
		Side:      BuySide,
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

	assert.Equal(t, Price(999), bb)
	assert.Equal(t, Price(1001), ba)
	assert.Equal(t, Price(1000), mid)
}

func printBook(_ *testing.T, lob *Orderbook) {
	for i := len(lob.asks.levels) - 1; i >= 0; i-- {
		ask := lob.asks.levels[i]
		fmt.Println("ask: ", "depth=", i, "size", ask.totalSize, "price=", ask.price)
	}

	for i, bid := range lob.bids.levels {
		fmt.Println("bid: ", "depth=", i, "size", bid.totalSize, "price=", bid.price)
	}

	fmt.Println()

}
