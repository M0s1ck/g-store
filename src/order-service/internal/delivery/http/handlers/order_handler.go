package handlers

import (
	"net/http"

	"orders-service/internal/delivery/http/dto"
	"orders-service/internal/delivery/http/helpers"
	"orders-service/internal/delivery/http/mapper"
	mymiddleware "orders-service/internal/delivery/http/middleware"
	"orders-service/internal/usecase/cancel_order"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/get_orders"
)

type OrderHandler struct {
	getById   *get_orders.GetByIdUsecase
	getByUser *get_orders.GetByUserUsecase
	create    *create_order.CreateOrderUsecase
	cancel    *cancel_order.CancelOrderUsecase
}

func NewOrderHandler(
	get *get_orders.GetByIdUsecase,
	getAll *get_orders.GetByUserUsecase,
	create *create_order.CreateOrderUsecase,
	cancel *cancel_order.CancelOrderUsecase,

) *OrderHandler {

	return &OrderHandler{
		getById:   get,
		getByUser: getAll,
		create:    create,
		cancel:    cancel,
	}
}

// GetById godoc
// @Summary Get order by id
// @Description Returns order by UUID
// @Tags orders_customers
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
		helpers.HandleError(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, mapper.OrderToResponse(order))
}

// GetByUser godoc
// @Summary Get all orders for user
// @Description Returns paginated list of orders for the authenticated user
// @Tags orders_customers
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
		helpers.HandleError(w, err)
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
// @Tags orders_customers
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)" example("123e4567-e89b-12d3-a456-426614174000")
// @Param order_request body dto.CreateOrderRequest true "Request to create an order"
// @Success 201 {object} dto.OrderCreatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := mymiddleware.UserIdFromContext(ctx)
	dtoReq := mymiddleware.BodyFromContext[dto.CreateOrderRequest](ctx)

	req := mapper.OrderCreateRequestToApplication(dtoReq)

	ordResp, err := h.create.Execute(ctx, req, userId)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	dtoResp := mapper.OrderCreatedResponseToDto(ordResp)
	helpers.RespondJSON(w, http.StatusCreated, dtoResp)
}

// Cancel godoc
// @Summary Cancels order (by customer)
// @Description Customer cancels order, body is optional. Possible reasons: CHANGED_MIND
// @Tags orders_customers
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Param X-User-ID header string true "User ID (UUID)" example("123e4567-e89b-12d3-a456-426614174000")
// @Param cancel_request body dto.ConsumerCancelOrderRequest false "cancel request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderID := mymiddleware.UUIDFromContext(ctx)
	userID := mymiddleware.UserIdFromContext(ctx)
	dtoReq := mymiddleware.OptionalBodyFromContext[dto.ConsumerCancelOrderRequest](ctx)

	req, err := mapper.ConsumerCancelOrderRequestToApplication(orderID, userID, dtoReq)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	err = h.cancel.Execute(ctx, req)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
