package dto

// Payment DTOs
type CreatePaymentRequest struct {
	ServiceID     string  `json:"service_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	PaymentDate   string  `json:"payment_date"`
	Notes         string  `json:"notes"`
}
