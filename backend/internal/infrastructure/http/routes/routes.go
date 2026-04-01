package routes

import (
	"mechanic-backend/internal/application/usecases"
	"mechanic-backend/internal/config"
	"mechanic-backend/internal/infrastructure/http/handlers"
	"mechanic-backend/internal/infrastructure/http/middleware"
	"mechanic-backend/internal/infrastructure/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Router struct {
	app *fiber.App
	cfg *config.Config
	db  *gorm.DB
}

func NewRouter(app *fiber.App, cfg *config.Config, db *gorm.DB) *Router {
	return &Router{
		app: app,
		cfg: cfg,
		db:  db,
	}
}

func (r *Router) Setup() {
	// Initialize repositories
	userRepo := repository.NewUserRepository(r.db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(r.db)
	vehicleRepo := repository.NewVehicleRepository(r.db)
	serviceRepo := repository.NewServiceRepository(r.db)
	paymentRepo := repository.NewPaymentRepository(r.db)

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(userRepo, refreshTokenRepo, r.cfg)
	vehicleUseCase := usecases.NewVehicleUseCase(vehicleRepo)
	serviceUseCase := usecases.NewServiceUseCase(serviceRepo, vehicleRepo)
	paymentUseCase := usecases.NewPaymentUseCase(paymentRepo, serviceRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	vehicleHandler := handlers.NewVehicleHandler(vehicleUseCase)
	serviceHandler := handlers.NewServiceHandler(serviceUseCase)
	paymentHandler := handlers.NewPaymentHandler(paymentUseCase)

	// Global middleware (applied to all routes)
	r.app.Use(middleware.Logger())
	r.app.Use(middleware.SecurityHeaders()) // Helmet equivalent: XSS, clickjacking protection
	r.app.Use(middleware.CORS(r.cfg))       // Restricted CORS to frontend URL
	r.app.Use(middleware.RateLimiter())     // General rate limiting: 100 req/min

	// Health check (public, no auth required)
	r.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "F&F Workshop API is running",
			"version": "1.0.0",
		})
	})

	// API v1 group
	api := r.app.Group("/api/v1")

	// Auth routes (public, with rate limiting)
	auth := api.Group("/auth")

	// Strict rate limiting for login/register (5 attempts/min to prevent brute force)
	auth.Post("/register", middleware.StrictRateLimiter(), authHandler.Register)
	auth.Post("/login", middleware.StrictRateLimiter(), authHandler.Login)

	// No rate limit for refresh (handled by single-flight pattern in frontend)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes (require authentication)
	protected := api.Group("", middleware.AuthMiddleware(r.cfg))

	// Logout routes (protected - require valid auth)
	protected.Post("/auth/logout", authHandler.Logout)               // New: logout (revoke refresh token)
	protected.Post("/auth/logout-all", authHandler.LogoutAllDevices) // New: logout from all devices

	// Vehicle routes
	vehicles := protected.Group("/vehicles")
	vehicles.Post("/", vehicleHandler.Create)
	vehicles.Get("/", vehicleHandler.GetAll)
	vehicles.Get("/:id", vehicleHandler.GetByID)
	vehicles.Put("/:id", vehicleHandler.Update)
	vehicles.Delete("/:id", vehicleHandler.Delete)

	// Service routes
	services := protected.Group("/services")
	services.Post("/quick", serviceHandler.CreateQuickTicket)
	services.Post("/", serviceHandler.Create)
	services.Get("/", serviceHandler.GetAll)
	services.Get("/:id", serviceHandler.GetByID)
	services.Put("/:id", serviceHandler.Update)
	services.Delete("/:id", serviceHandler.Delete)
	services.Get("/vehicle/:vehicleId", serviceHandler.GetByVehicleID)

	// Payment routes
	payments := protected.Group("/payments")
	payments.Post("/", paymentHandler.Create)
	payments.Get("/", paymentHandler.GetAll)
	payments.Get("/:id", paymentHandler.GetByID)
	payments.Delete("/:id", paymentHandler.Delete)
	payments.Get("/service/:serviceId", paymentHandler.GetByServiceID)

	// 404 handler
	r.app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
		})
	})
}
