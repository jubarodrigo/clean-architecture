package entity

import "context"

// OrderRepositoryInterface é a porta de persistência de pedidos.
type OrderRepositoryInterface interface {
	List(ctx context.Context) ([]*Order, error)
	Create(ctx context.Context, o *Order) error
}
