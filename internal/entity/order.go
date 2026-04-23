package entity

import "time"

// Order representa o pedido (entidade de domínio).
type Order struct {
	ID          string
	Customer    string
	Description string
	Amount      float64
	CreatedAt   time.Time
}
