package handlers

import (
	"errors"
	"net/http"
	"orders-service/internal/delivery/http/helpers"
	"orders-service/internal/delivery/http/mapper"
	derrors "orders-service/internal/domain/errors"
	"orders-service/internal/usecase/get_orders"

	mymiddleware "orders-service/internal/delivery/http/middleware"
)

type OrderHandler struct {
	get *get_orders.GetOrdersUsecase
}

func NewOrderHandler(get *get_orders.GetOrdersUsecase) *OrderHandler {
	return &OrderHandler{
		get: get,
	}
}

func (h *OrderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderId := mymiddleware.UUIDFromContext(r.Context())

	order, err := h.get.GetById(ctx, orderId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, mapper.ToResponse(order))
}

func (h *OrderHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, derrors.ErrOrderNotFound):
		helpers.RespondError(w, http.StatusNotFound, err.Error())

	default:
		helpers.RespondError(w, http.StatusInternalServerError, "internal error")
	}
}
