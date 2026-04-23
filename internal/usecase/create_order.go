package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/entity"
)

// CreateOrderInput representa os dados necessários para criar um pedido.
type CreateOrderInput struct {
	Customer    string
	Description string
	Amount      float64
}

// CreateOrderOutput contém o pedido criado.
type CreateOrderOutput struct {
	Order *entity.Order
}

// CreateOrderUseCase cria um novo pedido (útil para popular o banco nos testes).
type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

// NewCreateOrderUseCase construtor alinhado ao padrão do curso (wire).
func NewCreateOrderUseCase(r entity.OrderRepositoryInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{OrderRepository: r}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, in CreateOrderInput) (*CreateOrderOutput, error) {
	if uc.OrderRepository == nil {
		return nil, fmt.Errorf("create order: repositório não configurado")
	}
	in.Customer = strings.TrimSpace(in.Customer)
	in.Description = strings.TrimSpace(in.Description)
	if in.Customer == "" || in.Description == "" {
		return nil, fmt.Errorf("create order: cliente e descrição são obrigatórios")
	}
	if in.Amount <= 0 {
		return nil, fmt.Errorf("create order: valor deve ser positivo")
	}
	o := &entity.Order{
		Customer:    in.Customer,
		Description: in.Description,
		Amount:      in.Amount,
	}
	if err := uc.OrderRepository.Create(ctx, o); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	return &CreateOrderOutput{Order: o}, nil
}
