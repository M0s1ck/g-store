package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"orders-service/internal/delivery/http/dto"
	"orders-service/internal/delivery/http/helpers"
	"orders-service/internal/delivery/http/mapper"
	mymiddleware "orders-service/internal/delivery/http/middleware"
	derrors "orders-service/internal/domain/errors"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/get_orders"
)

type OrderHandler struct {
	getById   *get_orders.GetByIdUsecase
	getByUser *get_orders.GetByUserUsecase
	create    *create_order.CreateOrderUsecase
}

func NewOrderHandler(
	get *get_orders.GetByIdUsecase,
	getAll *get_orders.GetByUserUsecase,
	create *create_order.CreateOrderUsecase,
) *OrderHandler {

	return &OrderHandler{
		getById:   get,
		getByUser: getAll,
		create:    create,
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

	order, err := h.getById.Execute(ctx, orderId, userId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, mapper.OrderToResponse(order))
}

// GetByUser godoc
// @Summary Get all orders for user
// @Description Returns paginated list of orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.OrdersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := mymiddleware.UserIdFromContext(ctx)
	page := mymiddleware.PageFromContext(ctx)
	limit := mymiddleware.LimitFromContext(ctx)

	orders, total, err := h.getByUser.Execute(ctx, userId, page, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := dto.OrdersResponse{
		Orders: mapper.OrdersToResponse(orders),
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	helpers.RespondJSON(w, http.StatusOK, response)
}

// Create godoc
// @Summary Create an order
// @Description Creates a new order, a message is sent to payment service
// @Tags orders
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)"
// @Param order_request body dto.CreateOrderRequest true "Request to create an order"
// @Success 200 {object} dto.OrdersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	userId := mymiddleware.UserIdFromContext(ctx)

	req := mapper.OrderCreateRequestToApplication(dtoReq)

	ordResp, err := h.create.Execute(ctx, req, userId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	dtoResp := mapper.OrderCreatedResponseToDto(ordResp)
	helpers.RespondJSON(w, http.StatusCreated, dtoResp)
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
