package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"payment-service/internal/delivery/http/dto"
	"payment-service/internal/delivery/http/helpers"
	"payment-service/internal/delivery/http/mapper"
	mymiddleware "payment-service/internal/delivery/http/middleware"
	myerrors "payment-service/internal/domain/errors"
	"payment-service/internal/usecase/create_account"
	"payment-service/internal/usecase/get_account"
	"payment-service/internal/usecase/top_up"
)

type AccountHandler struct {
	getByID *get_account.GetByIdUsecase
	create  *create_account.CreateAccountUsecase
	topUp   *top_up.TopUpUsecase
}

func NewAccountHandler(
	getByID *get_account.GetByIdUsecase,
	create *create_account.CreateAccountUsecase,
	topUp *top_up.TopUpUsecase,
) *AccountHandler {
	return &AccountHandler{
		getByID: getByID,
		create:  create,
		topUp:   topUp,
	}
}

// GetById godoc
// @Summary Get account by id
// @Description Returns account by UUID
// @Tags accounts
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)" example("123e4567-e89b-12d3-a456-426614174000")
// @Param id path string true "Account ID (UUID)"
// @Success 200 {object} dto.AccountResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accId := mymiddleware.UUIDFromContext(ctx)
	userId := mymiddleware.UserIdFromContext(ctx)

	acc, err := h.getByID.Execute(ctx, accId, userId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	accResp := mapper.AccountToDtoResponse(acc)
	helpers.RespondJSON(w, http.StatusOK, accResp)
}

// Create godoc
// @Summary Create a new account
// @Description Creates a new account for user unless they already have one
// @Tags accounts
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)" example("123e4567-e89b-12d3-a456-426614174000")
// @Success 201 {object} dto.AccountCreatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := mymiddleware.UserIdFromContext(ctx)

	accResp, err := h.create.Execute(ctx, userId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	dtoResp := mapper.CreateResponseToDto(accResp)
	helpers.RespondJSON(w, http.StatusCreated, dtoResp)
}

// TopUp godoc
// @Summary Top up the account
// @Description Add given amount to account's balance (if it's user's)
// @Tags accounts
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID (UUID)" example("123e4567-e89b-12d3-a456-426614174000")
// @Param id path string true "Account ID (UUID)"
// @Param top_up_request body dto.TopUpReq true "Request to top up an account"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /accounts/{id} [patch]
func (h *AccountHandler) TopUp(w http.ResponseWriter, r *http.Request) {
	var dtoReq dto.TopUpReq
	if err := json.NewDecoder(r.Body).Decode(&dtoReq); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()

	accId := mymiddleware.UUIDFromContext(ctx)
	userId := mymiddleware.UserIdFromContext(ctx)

	err := h.topUp.Execute(ctx, accId, userId, dtoReq.Amount)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AccountHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, myerrors.ErrAccountNotFound):
		helpers.RespondError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, myerrors.ErrForbidden):
		helpers.RespondError(w, http.StatusForbidden, err.Error())

	case errors.Is(err, myerrors.ErrAmountNotPositive):
		helpers.RespondError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, myerrors.ErrAccountForUserAlreadyExists):
		helpers.RespondError(w, http.StatusConflict, err.Error())

	default:
		helpers.RespondError(w, http.StatusInternalServerError, "internal error: "+err.Error())
	}
}
