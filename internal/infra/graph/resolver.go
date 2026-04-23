package graph

import "github.com/rodrigocavalhero/clean_arch_orders_list/internal/usecase"

// Resolver concentra dependências para o GraphQL (gqlgen), no estilo 20-CleanArch.
type Resolver struct {
	ListOrdersUseCase  *usecase.ListOrdersUseCase
	CreateOrderUseCase *usecase.CreateOrderUseCase
}
