package usecase

import (
	"context"
	"fmt"

	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/entity"
)

// ListOrdersInput não possui filtros neste desafio.
type ListOrdersInput struct{}

// ListOrdersOutput agrupa o resultado do caso de uso.
type ListOrdersOutput struct {
	Orders []*entity.Order
}

// ListOrdersUseCase lista todos os pedidos.
type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

// NewListOrdersUseCase construtor alinhado ao padrão do curso (wire).
func NewListOrdersUseCase(r entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepository: r}
}

func (uc *ListOrdersUseCase) Execute(ctx context.Context, _ ListOrdersInput) (*ListOrdersOutput, error) {
	if uc.OrderRepository == nil {
		return nil, fmt.Errorf("list orders: repositório não configurado")
	}
	orders, err := uc.OrderRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	if orders == nil {
		orders = []*entity.Order{}
	}
	return &ListOrdersOutput{Orders: orders}, nil
}
