package service

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/rodrigocavalhero/clean_arch_orders_list/internal/infra/grpc/pb/order/v1"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/usecase"
)

// OrderService implementa o serviço gRPC (padrão internal/infra/grpc/service do curso).
type OrderService struct {
	orderv1.UnimplementedOrderServiceServer
	ListOrdersUseCase *usecase.ListOrdersUseCase
}

// NewOrderService construtor.
func NewOrderService(list *usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{ListOrdersUseCase: list}
}

// ListOrders delega para o caso de uso.
func (s *OrderService) ListOrders(ctx context.Context, _ *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	if s.ListOrdersUseCase == nil {
		return nil, status.Error(codes.Internal, "use case não configurado")
	}
	out, err := s.ListOrdersUseCase.Execute(ctx, usecase.ListOrdersInput{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orders := make([]*orderv1.Order, 0, len(out.Orders))
	for _, o := range out.Orders {
		orders = append(orders, &orderv1.Order{
			Id:          o.ID,
			Customer:    o.Customer,
			Description: o.Description,
			Amount:      o.Amount,
			CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		})
	}
	return &orderv1.ListOrdersResponse{Orders: orders}, nil
}
