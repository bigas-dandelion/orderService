package handler

import (
	"encoding/json"
	"l0/cons/internal/services"
	"net/http"
)

type HandlerTask struct {
	service *services.Service
}

func NewHandlerTask(service *services.Service) *HandlerTask {
	return &HandlerTask{
		service: service,
	}
}

func (h *HandlerTask) GetOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("order_uid")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		order, err := h.service.GetOrder(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if order == nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)
	}
}
