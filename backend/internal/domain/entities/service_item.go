package entities

import (
	"time"

	"github.com/google/uuid"
)

type ServiceItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ServiceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Description string    `gorm:"type:varchar(255);not null" json:"description"`
	Quantity    int       `gorm:"type:integer;not null;default:1" json:"quantity"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
}

func (ServiceItem) TableName() string {
	return "service_items"
}
