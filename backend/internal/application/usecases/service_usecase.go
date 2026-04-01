package usecases

import (
	"mechanic-backend/internal/application/dto"
	"mechanic-backend/internal/domain/entities"
	"mechanic-backend/internal/domain/repositories"

	"github.com/google/uuid"
)

type ServiceUseCase struct {
	serviceRepo repositories.ServiceRepository
	vehicleRepo repositories.VehicleRepository
}

func NewServiceUseCase(serviceRepo repositories.ServiceRepository, vehicleRepo repositories.VehicleRepository) *ServiceUseCase {
	return &ServiceUseCase{
		serviceRepo: serviceRepo,
		vehicleRepo: vehicleRepo,
	}
}

func (uc *ServiceUseCase) Create(req dto.CreateServiceRequest) (*entities.Service, error) {
	vehicleID, err := uuid.Parse(req.VehicleID)
	if err != nil {
		return nil, err
	}

	// Verify vehicle exists
	_, err = uc.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return nil, err
	}

	// Calculate total cost from items
	var totalCost float64
	for _, item := range req.Items {
		totalCost += float64(item.Quantity) * item.Price
	}

	service := &entities.Service{
		VehicleID:   vehicleID,
		Description: req.Description,
		Cost:        totalCost,
		Status:      entities.StatusPending,
	}

	if err := uc.serviceRepo.Create(service); err != nil {
		return nil, err
	}

	// Create service items
	for _, itemReq := range req.Items {
		item := &entities.ServiceItem{
			ServiceID:   service.ID,
			Description: itemReq.Description,
			Quantity:    itemReq.Quantity,
			Price:       itemReq.Price,
		}
		if err := uc.serviceRepo.CreateItem(item); err != nil {
			return nil, err
		}
	}

	// Load items for response
	service, _ = uc.serviceRepo.FindByIDWithItems(service.ID)
	return service, nil
}

func (uc *ServiceUseCase) GetByID(id uuid.UUID) (*entities.Service, error) {
	return uc.serviceRepo.FindByID(id)
}

func (uc *ServiceUseCase) GetByVehicleID(vehicleID uuid.UUID) ([]entities.Service, error) {
	return uc.serviceRepo.FindByVehicleID(vehicleID)
}

func (uc *ServiceUseCase) GetAll() ([]entities.Service, error) {
	return uc.serviceRepo.FindAll()
}

func (uc *ServiceUseCase) GetAllWithPagination(page, limit int) ([]entities.Service, int64, error) {
	return uc.serviceRepo.FindAllWithPagination(page, limit)
}

func (uc *ServiceUseCase) Update(id uuid.UUID, req dto.UpdateServiceRequest) (*entities.Service, error) {
	service, err := uc.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Description != "" {
		service.Description = req.Description
	}

	// Update items if provided
	if len(req.Items) > 0 {
		// Delete old items
		if err := uc.serviceRepo.DeleteItemsByServiceID(id); err != nil {
			return nil, err
		}

		// Create new items and calculate cost
		var totalCost float64
		for _, itemReq := range req.Items {
			totalCost += float64(itemReq.Quantity) * itemReq.Price
			item := &entities.ServiceItem{
				ServiceID:   id,
				Description: itemReq.Description,
				Quantity:    itemReq.Quantity,
				Price:       itemReq.Price,
			}
			if err := uc.serviceRepo.CreateItem(item); err != nil {
				return nil, err
			}
		}
		service.Cost = totalCost
	}

	if req.Status != "" {
		service.Status = entities.ServiceStatus(req.Status)
	}

	if err := uc.serviceRepo.Update(service); err != nil {
		return nil, err
	}

	// Load items for response
	service, _ = uc.serviceRepo.FindByIDWithItems(service.ID)
	return service, nil
}

func (uc *ServiceUseCase) Delete(id uuid.UUID) error {
	return uc.serviceRepo.Delete(id)
}

func (uc *ServiceUseCase) CreateQuickTicket(req dto.QuickTicketRequest) (*entities.Service, error) {
	vehicleModel := req.VehiculoDesc
	if req.VehiculoColor != "" {
		vehicleModel += " - " + req.VehiculoColor
	}

	vehicle := &entities.Vehicle{
		ClientName:   req.ClienteNombre,
		Phone:        req.ClienteTelefono,
		Brand:        req.VehiculoDesc,  // Store description in Brand so it shows up in UI
		Model:        req.VehiculoColor, // Store color in Model (can be empty string)
		VehicleModel: vehicleModel,
		PlateNumber:  "N/A",
	}

	if err := uc.vehicleRepo.Create(vehicle); err != nil {
		return nil, err
	}

	var totalCost float64
	for _, srv := range req.Servicios {
		quantity := srv.Cantidad
		if quantity <= 0 {
			quantity = 1
		}
		totalCost += srv.Costo * float64(quantity)
	}

	service := &entities.Service{
		VehicleID:   vehicle.ID,
		Description: "Recibo rapido",
		Cost:        totalCost,
		Status:      entities.StatusPending,
	}

	if err := uc.serviceRepo.Create(service); err != nil {
		return nil, err
	}

	for _, srv := range req.Servicios {
		quantity := srv.Cantidad
		if quantity <= 0 {
			quantity = 1
		}
		item := &entities.ServiceItem{
			ServiceID:   service.ID,
			Description: srv.Descripcion,
			Quantity:    quantity,
			Price:       srv.Costo,
		}
		if err := uc.serviceRepo.CreateItem(item); err != nil {
			return nil, err
		}
	}

	service, _ = uc.serviceRepo.FindByIDWithItems(service.ID)
	service.Vehicle = vehicle
	return service, nil
}
