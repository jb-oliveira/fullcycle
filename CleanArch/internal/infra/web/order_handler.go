package web

import (
	"encoding/json"
	"net/http"

	"github.com/jb-oliveira/fullcycle/CleanArch/internal/usecase"
)

type OrderHandler struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
}

func NewOrderHandler(createOrderUseCase *usecase.CreateOrderUseCase) *OrderHandler {
	return &OrderHandler{CreateOrderUseCase: createOrderUseCase}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.CreateOrderUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
