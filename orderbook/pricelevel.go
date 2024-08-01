package orderbook

type (
	Size  float64
	Price float64
)

func NewPriceLevel(price Price) *PriceLevel {
	return &PriceLevel{
		orders: make([]*Order, 0),
	}
}

type PriceLevel struct {
	orders    []*Order
	price     Price
	totalSize Size
}
