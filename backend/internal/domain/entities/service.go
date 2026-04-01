package entities

import (
	"time"

	"github.com/google/uuid"
)

type ServiceStatus string

const (
	StatusPending    ServiceStatus = "pending"
	StatusInProgress ServiceStatus = "in_progress"
	StatusCompleted  ServiceStatus = "completed"
	StatusCancelled  ServiceStatus = "cancelled"
)

type Service struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VehicleID   uuid.UUID     `gorm:"type:uuid;not null;index" json:"vehicle_id"`
	Vehicle     *Vehicle      `gorm:"foreignKey:VehicleID;constraint:OnDelete:CASCADE" json:"vehicle,omitempty"`
	Description string        `gorm:"type:text" json:"description"`
	Cost        float64       `gorm:"type:decimal(10,2);not null;default:0.00" json:"cost"`
	Status      ServiceStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	CreatedAt   time.Time     `gorm:"default:now();index" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"default:now()" json:"updated_at"`
	Items       []ServiceItem `gorm:"foreignKey:ServiceID" json:"items,omitempty"`
	Payments    []Payment     `gorm:"foreignKey:ServiceID" json:"payments,omitempty"`
}

func (Service) TableName() string {
	return "services"
}
