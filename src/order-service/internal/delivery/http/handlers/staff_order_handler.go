package handlers

import (
	"net/http"

	"orders-service/internal/delivery/http/dto"
	"orders-service/internal/delivery/http/helpers"
	"orders-service/internal/delivery/http/mapper"
	mymiddleware "orders-service/internal/delivery/http/middleware"
	"orders-service/internal/usecase/cancel_order"
	"orders-service/internal/usecase/order_update_status"
)

type StaffOrderHandler struct {
	cancelUC       *cancel_order.CancelOrderUsecase
	updateStatusUC *order_update_status.UpdateStatusUsecase
}

func NewStaffOrderHandler(
	cancelUC *cancel_order.CancelOrderUsecase,
	updateStatusUC *order_update_status.UpdateStatusUsecase,
) *StaffOrderHandler {

	return &StaffOrderHandler{
		cancelUC:       cancelUC,
		updateStatusUC: updateStatusUC,
	}
}

// Cancel godoc
// @Summary Cancels order (by store)
// @Description Cancels order, used by store/staff, simulation of jwt role with api key that can be found in .env
// @Description Possible reasons: OUT_OF_STOCK, DELIVERY_UNAVAILABLE
// @Tags orders_staff
// @Accept json
// @Produce json
// @Param X-Staff-API-Key header string true "Staff API key"
// @Param id path string true "Order ID (UUID)"
// @Param cancel_request body dto.StaffCancelOrderRequest true "cancel request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /staff/orders/{id}/cancel [post]
func (h *StaffOrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderID := mymiddleware.UUIDFromContext(ctx)
	dtoReq := mymiddleware.BodyFromContext[dto.StaffCancelOrderRequest](ctx)

	req, err := mapper.StaffCancelOrderRequestToApplication(orderID, dtoReq)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	err = h.cancelUC.Execute(ctx, req)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateStatus godoc
// @Summary Update order status (by store/stuff)
// @Description Updates order status, used by staff, simulation of jwt role with api key that can be found in .env
// @Description Possible update pipeline: PAID -> ASSEMBLING -> ASSEMBLED -> DELIVERING -> ISSUED
// @Tags orders_staff
// @Accept json
// @Produce json
// @Param X-Staff-API-Key header string true "Staff API key"
// @Param id path string true "Order ID (UUID)"
// @Param status_request body dto.UpdateOrderStatusRequest true "Order status update request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /staff/orders/{id}/status [patch]
func (h *StaffOrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderID := mymiddleware.UUIDFromContext(ctx)
	dtoReq := mymiddleware.BodyFromContext[dto.UpdateOrderStatusRequest](ctx)

	req, err := mapper.UpdateOrderStatusRequestToApplication(orderID, dtoReq)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	err = h.updateStatusUC.Execute(ctx, req)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
