package repository

import (
	"mechanic-backend/internal/domain/entities"
	"mechanic-backend/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) repositories.ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(service *entities.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) FindByID(id uuid.UUID) (*entities.Service, error) {
	var service entities.Service
	err := r.db.Preload("Vehicle").Preload("Payments").First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) FindByIDWithItems(id uuid.UUID) (*entities.Service, error) {
	var service entities.Service
	err := r.db.Preload("Vehicle").Preload("Items").Preload("Payments").First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) FindByVehicleID(vehicleID uuid.UUID) ([]entities.Service, error) {
	var services []entities.Service
	err := r.db.Preload("Payments").Find(&services, "vehicle_id = ?", vehicleID).Error
	return services, err
}

func (r *serviceRepository) FindAll() ([]entities.Service, error) {
	var services []entities.Service
	err := r.db.Preload("Vehicle").Preload("Items").Order("created_at DESC").Find(&services).Error
	return services, err
}

func (r *serviceRepository) FindAllWithPagination(page, limit int) ([]entities.Service, int64, error) {
	var services []entities.Service
	var total int64

	// Count total records
	if err := r.db.Model(&entities.Service{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated records ordered by created_at DESC (newest first)
	err := r.db.Preload("Vehicle").Preload("Items").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&services).Error

	return services, total, err
}

func (r *serviceRepository) Update(service *entities.Service) error {
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Service{}, "id = ?", id).Error
}

func (r *serviceRepository) CreateItem(item *entities.ServiceItem) error {
	return r.db.Create(item).Error
}

func (r *serviceRepository) DeleteItemsByServiceID(serviceID uuid.UUID) error {
	return r.db.Delete(&entities.ServiceItem{}, "service_id = ?", serviceID).Error
}
