package repositories

import (
	"mechanic-backend/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByID(id uuid.UUID) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindAll() ([]entities.User, error)
	Update(user *entities.User) error
	Delete(id uuid.UUID) error
}
