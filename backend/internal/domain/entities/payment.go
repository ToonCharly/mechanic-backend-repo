package entities

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "cash"
	PaymentCard     PaymentMethod = "card"
	PaymentTransfer PaymentMethod = "transfer"
	PaymentOther    PaymentMethod = "other"
)

type Payment struct {
	ID            uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ServiceID     uuid.UUID     `gorm:"type:uuid;not null;index" json:"service_id"`
	Service       *Service      `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	Amount        float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentMethod PaymentMethod `gorm:"type:varchar(20);not null;default:'cash';index" json:"payment_method"`
	PaymentDate   time.Time     `gorm:"type:date;not null;default:current_date;index" json:"payment_date"`
	Notes         string        `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time     `gorm:"default:now()" json:"created_at"`
}

func (Payment) TableName() string {
	return "payments"
}
