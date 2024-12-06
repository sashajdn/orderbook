package client

import (
	"context"

	"github.com/sashajdn/orderbook/lob"
)

type AddOrderRequest struct {
	OrderType lob.OrderType
	OrderSide lob.OrderSide
	Price     lob.Price
	Size      lob.Size
}

type AddOrderResponse struct {
	OrderID uint64
}

type CancelOrderResponse struct{}
type CancelOrderRequest struct{}

type EditOrderRequest struct{}
type EditOrderResponse struct{}

type Client interface {
	AddOrder(ctx context.Context, req AddOrderRequest) (AddOrderResponse, error)
	CancelOrder(ctx context.Context, req CancelOrderRequest) (CancelOrderResponse, error)
	EditOrder(ctx context.Context, req EditOrderRequest) (EditOrderResponse, error)
}
