package usecases

import (
	"mechanic-backend/internal/application/dto"
	"mechanic-backend/internal/domain/entities"
	"mechanic-backend/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	paymentRepo repositories.PaymentRepository
	serviceRepo repositories.ServiceRepository
}

func NewPaymentUseCase(paymentRepo repositories.PaymentRepository, serviceRepo repositories.ServiceRepository) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo: paymentRepo,
		serviceRepo: serviceRepo,
	}
}

func (uc *PaymentUseCase) Create(req dto.CreatePaymentRequest) (*entities.Payment, error) {
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}

	// Verify service exists
	_, err = uc.serviceRepo.FindByID(serviceID)
	if err != nil {
		return nil, err
	}

	paymentDate := time.Now()
	if req.PaymentDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.PaymentDate)
		if err == nil {
			paymentDate = parsedDate
		}
	}

	payment := &entities.Payment{
		ServiceID:     serviceID,
		Amount:        req.Amount,
		PaymentMethod: entities.PaymentMethod(req.PaymentMethod),
		PaymentDate:   paymentDate,
		Notes:         req.Notes,
	}

	if err := uc.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *PaymentUseCase) GetByID(id uuid.UUID) (*entities.Payment, error) {
	return uc.paymentRepo.FindByID(id)
}

func (uc *PaymentUseCase) GetByServiceID(serviceID uuid.UUID) ([]entities.Payment, error) {
	return uc.paymentRepo.FindByServiceID(serviceID)
}

func (uc *PaymentUseCase) GetAll() ([]entities.Payment, error) {
	return uc.paymentRepo.FindAll()
}

func (uc *PaymentUseCase) Delete(id uuid.UUID) error {
	return uc.paymentRepo.Delete(id)
}
