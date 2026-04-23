//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/entity"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/infra/database"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/infra/web"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/usecase"
)

var setRepo = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

// NewListOrdersUseCase expõe o wire injector (padrão do curso).
func NewListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(setRepo, usecase.NewListOrdersUseCase)
	return nil
}

// NewCreateOrderUseCase expõe o wire injector.
func NewCreateOrderUseCase(db *sql.DB) *usecase.CreateOrderUseCase {
	wire.Build(setRepo, usecase.NewCreateOrderUseCase)
	return nil
}

// NewWebOrderHandler compõe os casos de uso via wire.
func NewWebOrderHandler(db *sql.DB) *web.WebOrderHandler {
	wire.Build(setRepo, usecase.NewListOrdersUseCase, usecase.NewCreateOrderUseCase, web.NewWebOrderHandler)
	return nil
}
