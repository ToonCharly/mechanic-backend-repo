package entities

import (
	"time"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ClientName   string    `gorm:"type:varchar(100);not null;index" json:"client_name"`
	Phone        string    `gorm:"type:varchar(20)" json:"phone"`
	VehicleModel string    `gorm:"type:varchar(50)" json:"vehicle_model"`
	Brand        string    `gorm:"type:varchar(50)" json:"brand"`
	Model        string    `gorm:"type:varchar(50)" json:"model"`
	Year         int       `gorm:"type:integer" json:"year"`
	Color        string    `gorm:"type:varchar(30)" json:"color"`
	PlateNumber  string    `gorm:"type:varchar(20);index" json:"plate_number"`
	CreatedAt    time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:now()" json:"updated_at"`
	Services     []Service `gorm:"foreignKey:VehicleID;constraint:OnDelete:CASCADE" json:"services,omitempty"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}
