package lob

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBook_SanityChecks(t *testing.T) {
	fmt.Println(`========= Buy Side ===========`)
	buySide := NewBook(BuySide)
	buyOrders := generateOrders(10, 0, 990, 1, []uint64{1, 2, 4, 32})

	for _, order := range buyOrders {
		buySide.Make(order)
	}

	fmt.Println("buy depth: ", buySide.Depth())

	fills, err := buySide.Take(6.5)
	require.NoError(t, err)

	fmt.Println("buy depth: ", buySide.Depth())

	for _, fill := range fills {
		fmt.Println("Buy Side Fill", fill.String())
	}

	assert.True(t, len(fills) >= 1)
}
