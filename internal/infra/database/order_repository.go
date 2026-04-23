package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/entity"
)

// OrderRepository implementa entity.OrderRepositoryInterface com PostgreSQL (via database/sql + pgx).
type OrderRepository struct {
	DB *sql.DB
}

// NewOrderRepository segue o padrão de construtor do curso (injeção de *sql.DB).
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

// List retorna pedidos ordenados do mais recente.
func (r *OrderRepository) List(ctx context.Context) ([]*entity.Order, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("order repository: db nulo")
	}
	const q = `
		SELECT id, customer_name, description, amount, created_at
		FROM orders
		ORDER BY created_at DESC, id
	`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.Customer, &o.Description, &o.Amount, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = []*entity.Order{}
	}
	return out, nil
}

// Create insere e preenche id e created_at no ponteiro.
func (r *OrderRepository) Create(ctx context.Context, o *entity.Order) error {
	if r.DB == nil {
		return fmt.Errorf("order repository: db nulo")
	}
	const q = `
		INSERT INTO orders (customer_name, description, amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.DB.QueryRowContext(ctx, q, o.Customer, o.Description, o.Amount).Scan(&o.ID, &o.CreatedAt)
}
