package repositories

import (
	"mechanic-backend/internal/domain/entities"

	"github.com/google/uuid"
)

type VehicleRepository interface {
	Create(vehicle *entities.Vehicle) error
	FindByID(id uuid.UUID) (*entities.Vehicle, error)
	FindByPlateNumber(plateNumber string) (*entities.Vehicle, error)
	FindAll() ([]entities.Vehicle, error)
	Update(vehicle *entities.Vehicle) error
	Delete(id uuid.UUID) error
}
