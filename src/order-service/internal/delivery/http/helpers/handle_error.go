package helpers

import (
	"errors"
	"net/http"
	derrors "orders-service/internal/domain/errors"
)

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, derrors.ErrOrderNotFound):
		RespondError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, derrors.ErrForbidden):
		RespondError(w, http.StatusForbidden, err.Error())

	case errors.Is(err, derrors.ErrAmountNotPositive):
		RespondError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, derrors.ErrInvalidOrderStatus):
		RespondError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, derrors.ErrCancellationReasonRequired):
		RespondError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, derrors.ErrInvalidOrderStatusChange):
		RespondError(w, http.StatusConflict, err.Error())

	case errors.Is(err, derrors.ErrInvalidCancellationReason):
		RespondError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, derrors.ErrOrderAlreadyCanceled):
		RespondError(w, http.StatusConflict, err.Error())

	case errors.Is(err, derrors.ErrOrderCantBeCancelled):
		RespondError(w, http.StatusConflict, err.Error())

	default:
		RespondError(w, http.StatusInternalServerError, "internal error: "+err.Error())
	}
}
