package repositories

import (
	"mechanic-backend/internal/domain/entities"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(payment *entities.Payment) error
	FindByID(id uuid.UUID) (*entities.Payment, error)
	FindByServiceID(serviceID uuid.UUID) ([]entities.Payment, error)
	FindAll() ([]entities.Payment, error)
	Delete(id uuid.UUID) error
}
