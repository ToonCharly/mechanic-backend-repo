package database

import (
	"log"
	"mechanic-backend/internal/config"
	"mechanic-backend/internal/domain/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var err error

	logLevel := logger.Info
	if cfg.App.Environment == "production" {
		logLevel = logger.Error
	}

	DB, err = gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	log.Println("✅ Database connected successfully")

	// Auto migrate tables
	if err := AutoMigrate(DB); err != nil {
		return nil, err
	}

	return DB, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	// Migrate in order: parent tables first, then child tables
	tables := []interface{}{
		&entities.User{},
		&entities.RefreshToken{}, // Session management
		&entities.Vehicle{},
		&entities.Service{},
		&entities.ServiceItem{}, // Service line items
		&entities.Payment{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			log.Printf("❌ Migration failed for %T: %v", table, err)
			return err
		}
	}

	log.Println("✅ Database migrations completed")

	// Drop unique constraint on plate_number to allow vehicles without plates
	db.Exec("ALTER TABLE vehicles ALTER COLUMN plate_number DROP NOT NULL")
	db.Exec("DROP INDEX IF EXISTS idx_vehicles_plate_number")

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
