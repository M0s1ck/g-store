package mapper

import (
	"payment-service/internal/delivery/http/dto"
	"payment-service/internal/domain/entities"
	"payment-service/internal/usecase/create_account"
)

func CreateResponseToDto(response *create_account.Response) *dto.AccountCreatedResponse {
	return &dto.AccountCreatedResponse{
		AccountID: response.ID,
	}
}

func AccountToDtoResponse(acc *entities.Account) *dto.AccountResponse {
	return &dto.AccountResponse{
		ID:        acc.ID,
		UserID:    acc.UserID,
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
}
