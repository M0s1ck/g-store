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

// GetById godoc
// @Summary Get order by id
// @Description Returns order by UUID
// @Tags orders
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderId := mymiddleware.UUIDFromContext(ctx)
	userId := mymiddleware.UserIdFromContext(ctx)

	order, err := h.get.GetById(ctx, orderId, userId)
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

	case errors.Is(err, derrors.ErrForbidden):
		helpers.RespondError(w, http.StatusForbidden, err.Error())

	default:
		helpers.RespondError(w, http.StatusInternalServerError, "internal error: "+err.Error())
	}
}
