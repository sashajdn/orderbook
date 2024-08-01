package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBook(t *testing.T) {
	t.Parallel()

	book := make(BookLevels, 128)

	book.Make(&Order{Price: 1.0, Size: 1.0})
	book.Make(&Order{Price: 2.0, Size: 1.0})
	book.Make(&Order{Price: 3.0, Size: 1.0})
	book.Make(&Order{Price: 4.0, Size: 1.0})
	book.Make(&Order{Price: 5.0, Size: 1.0})
	book.Make(&Order{Price: 1.0, Size: 1.0})

	assert.Equal(t, 5, book.Depth())

	top, err := book.Top()
	require.NoError(t, err)
	assert.Equal(t, 1.0, top)
}
