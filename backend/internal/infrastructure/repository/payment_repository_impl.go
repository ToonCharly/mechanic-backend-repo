package repository

import (
	"mechanic-backend/internal/domain/entities"
	"mechanic-backend/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) repositories.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *entities.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.Preload("Service").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByServiceID(serviceID uuid.UUID) ([]entities.Payment, error) {
	var payments []entities.Payment
	err := r.db.Find(&payments, "service_id = ?", serviceID).Error
	return payments, err
}

func (r *paymentRepository) FindAll() ([]entities.Payment, error) {
	var payments []entities.Payment
	err := r.db.Preload("Service").Preload("Service.Vehicle").Preload("Service.Items").Order("created_at desc").Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Payment{}, "id = ?", id).Error
}
