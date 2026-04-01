package usecases

import (
	"mechanic-backend/internal/application/dto"
	"mechanic-backend/internal/domain/entities"
	"mechanic-backend/internal/domain/repositories"

	"github.com/google/uuid"
)

type VehicleUseCase struct {
	vehicleRepo repositories.VehicleRepository
}

func NewVehicleUseCase(vehicleRepo repositories.VehicleRepository) *VehicleUseCase {
	return &VehicleUseCase{vehicleRepo: vehicleRepo}
}

func (uc *VehicleUseCase) Create(req dto.CreateVehicleRequest) (*entities.Vehicle, error) {
	vehicle := &entities.Vehicle{
		ClientName:  req.ClientName,
		Phone:       req.Phone,
		Brand:       req.Brand,
		Model:       req.Model,
		Year:        req.Year,
		Color:       req.Color,
		PlateNumber: req.PlateNumber,
	}

	if err := uc.vehicleRepo.Create(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (uc *VehicleUseCase) GetByID(id uuid.UUID) (*entities.Vehicle, error) {
	return uc.vehicleRepo.FindByID(id)
}

func (uc *VehicleUseCase) GetByPlateNumber(plateNumber string) (*entities.Vehicle, error) {
	return uc.vehicleRepo.FindByPlateNumber(plateNumber)
}

func (uc *VehicleUseCase) GetAll() ([]entities.Vehicle, error) {
	return uc.vehicleRepo.FindAll()
}

func (uc *VehicleUseCase) Update(id uuid.UUID, req dto.UpdateVehicleRequest) (*entities.Vehicle, error) {
	vehicle, err := uc.vehicleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ClientName != "" {
		vehicle.ClientName = req.ClientName
	}
	if req.Phone != "" {
		vehicle.Phone = req.Phone
	}
	if req.Brand != "" {
		vehicle.Brand = req.Brand
	}
	if req.Model != "" {
		vehicle.Model = req.Model
	}
	if req.Year != 0 {
		vehicle.Year = req.Year
	}
	if req.Color != "" {
		vehicle.Color = req.Color
	}
	if req.PlateNumber != "" {
		vehicle.PlateNumber = req.PlateNumber
	}

	if err := uc.vehicleRepo.Update(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (uc *VehicleUseCase) Delete(id uuid.UUID) error {
	return uc.vehicleRepo.Delete(id)
}
