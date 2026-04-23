package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/entity"
	"github.com/rodrigocavalhero/clean_arch_orders_list/internal/usecase"
)

const jsonContentType = "application/json; charset=utf-8"

// WebOrderHandler expõe os casos de uso via HTTP (REST), no estilo do curso.
type WebOrderHandler struct {
	ListOrdersUseCase  *usecase.ListOrdersUseCase
	CreateOrderUseCase *usecase.CreateOrderUseCase
}

// NewWebOrderHandler constrói o handler.
func NewWebOrderHandler(
	list *usecase.ListOrdersUseCase,
	create *usecase.CreateOrderUseCase,
) *WebOrderHandler {
	return &WebOrderHandler{ListOrdersUseCase: list, CreateOrderUseCase: create}
}

type orderDTO struct {
	ID          string  `json:"id"`
	Customer    string  `json:"customer"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	CreatedAt   string  `json:"created_at"`
}

type listOrdersResponse struct {
	Data []orderDTO `json:"data"`
}

type createOrderRequest struct {
	Customer    string  `json:"customer"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type createOrderResponse struct {
	Data orderDTO `json:"data"`
}

func toDTOs(orders []*entity.Order) []orderDTO {
	if orders == nil {
		return []orderDTO{}
	}
	out := make([]orderDTO, 0, len(orders))
	for _, o := range orders {
		out = append(out, orderDTO{
			ID:          o.ID,
			Customer:    o.Customer,
			Description: o.Description,
			Amount:      o.Amount,
			CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		})
	}
	return out
}

// List trata GET /order.
func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.ListOrdersUseCase.Execute(r.Context(), usecase.ListOrdersInput{})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", jsonContentType)
	_ = json.NewEncoder(w).Encode(listOrdersResponse{Data: toDTOs(out.Orders)})
}

// Create trata POST /order.
func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "body JSON inválido")
		return
	}
	out, err := h.CreateOrderUseCase.Execute(r.Context(), usecase.CreateOrderInput{
		Customer:    body.Customer,
		Description: body.Description,
		Amount:      body.Amount,
	})
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	o := out.Order
	_ = json.NewEncoder(w).Encode(createOrderResponse{Data: orderDTO{
		ID:          o.ID,
		Customer:    o.Customer,
		Description: o.Description,
		Amount:      o.Amount,
		CreatedAt:   o.CreatedAt.Format(time.RFC3339),
	}})
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
