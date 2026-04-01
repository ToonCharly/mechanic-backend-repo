package repository

import (
	"mechanic-backend/internal/domain/entities"
	"mechanic-backend/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) repositories.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(vehicle *entities.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *vehicleRepository) FindByID(id uuid.UUID) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle
	err := r.db.Preload("Services").First(&vehicle, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *vehicleRepository) FindByPlateNumber(plateNumber string) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle
	err := r.db.First(&vehicle, "plate_number = ?", plateNumber).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *vehicleRepository) FindAll() ([]entities.Vehicle, error) {
	var vehicles []entities.Vehicle
	err := r.db.Find(&vehicles).Error
	return vehicles, err
}

func (r *vehicleRepository) Update(vehicle *entities.Vehicle) error {
	return r.db.Save(vehicle).Error
}

func (r *vehicleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Vehicle{}, "id = ?", id).Error
}
