package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type AppConfig struct {
	Environment string
	Port        string
	FrontendURL string
}

func Load() *Config {
	// Load .env file (ignore error in production with env vars)
	_ = godotenv.Load()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "admin123")
	dbName := getEnv("DB_NAME", "mechanic_db")

	dsn := getEnv("DATABASE_URL", "")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName,
		)
	}

	jwtExpire, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "72"))
	if err != nil {
		log.Println("Invalid JWT_EXPIRE_HOURS, using default: 72")
		jwtExpire = 72
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			DSN:      dsn,
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-change-me"),
			ExpireHours: jwtExpire,
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("PORT", "8080"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
