package load

import (
	"context"
	"fmt"

	"github.com/sashajdn/orderbook/lob"
)

func NewLOBClient(lob *lob.Orderbook) *LOBClient {
	return &LOBClient{
		lob: lob,
	}
}

var _ Client = &LOBClient{}

type LOBClient struct {
	lob *lob.Orderbook
}

func (l *LOBClient) AddOrder(ctx context.Context, req AddOrderRequest) (AddOrderResponse, error) {
	order := &lob.Order{
		OrderType: req.OrderType,
		Side:      req.OrderSide,
		Price:     req.Price,
		Size:      req.Size,
	}

	id, err := l.lob.PlaceOrder(order)
	if err != nil {
		return AddOrderResponse{}, fmt.Errorf("add order: %w", err)
	}

	return AddOrderResponse{
		OrderID: id,
	}, nil
}

func (l *LOBClient) CancelOrder(ctx context.Context, req CancelOrderRequest) (CancelOrderResponse, error) {
	return CancelOrderResponse{}, fmt.Errorf(`unimplemented`)
}
func (l *LOBClient) EditOrder(ctx context.Context, req EditOrderRequest) (EditOrderResponse, error) {
	return EditOrderResponse{}, fmt.Errorf(`unimplemented`)
}
