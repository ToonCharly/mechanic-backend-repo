package main

import (
	"fmt"
	"mechanic-backend/internal/config"
	"mechanic-backend/internal/infrastructure/database"
	"mechanic-backend/internal/infrastructure/http/routes"
	"mechanic-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize structured logger
	if err := logger.Initialize(cfg.App.Environment); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting MechanicPro API",
		zap.String("environment", cfg.App.Environment),
		zap.String("port", cfg.App.Port),
	)

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "MechanicPro API v1.0",
		ServerHeader: "Fiber",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Log internal errors but don't expose them
			if code >= 500 {
				logger.Error("Internal server error",
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.Error(err),
				)
				return c.Status(code).JSON(fiber.Map{
					"success": false,
					"error":   "Internal server error",
				})
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Setup routes
	router := routes.NewRouter(app, cfg, db)
	router.Setup()

	// Start server
	port := fmt.Sprintf(":%s", cfg.App.Port)
	logger.Info("Server starting",
		zap.String("port", cfg.App.Port),
		zap.String("url", fmt.Sprintf("http://localhost%s", port)),
	)

	if err := app.Listen(port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
