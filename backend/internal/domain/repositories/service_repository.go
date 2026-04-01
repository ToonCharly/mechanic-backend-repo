package repositories

import (
	"mechanic-backend/internal/domain/entities"

	"github.com/google/uuid"
)

type ServiceRepository interface {
	Create(service *entities.Service) error
	FindByID(id uuid.UUID) (*entities.Service, error)
	FindByIDWithItems(id uuid.UUID) (*entities.Service, error)
	FindByVehicleID(vehicleID uuid.UUID) ([]entities.Service, error)
	FindAll() ([]entities.Service, error)
	FindAllWithPagination(page, limit int) ([]entities.Service, int64, error)
	Update(service *entities.Service) error
	Delete(id uuid.UUID) error
	CreateItem(item *entities.ServiceItem) error
	DeleteItemsByServiceID(serviceID uuid.UUID) error
}
