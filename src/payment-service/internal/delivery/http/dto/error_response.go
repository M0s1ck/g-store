package dto

type ErrorResponse struct {
	Error string `json:"error" example:"account not found: id=550e8400-e29b-41d4-a716-446655440000"`
}
