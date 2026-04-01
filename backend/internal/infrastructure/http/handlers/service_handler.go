package handlers

import (
	"mechanic-backend/internal/application/dto"
	"mechanic-backend/internal/application/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	serviceUseCase *usecases.ServiceUseCase
}

func NewServiceHandler(serviceUseCase *usecases.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{serviceUseCase: serviceUseCase}
}

func (h *ServiceHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	service, err := h.serviceUseCase.Create(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Service created successfully",
		"data":    service,
	})
}

func (h *ServiceHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid service ID",
		})
	}

	service, err := h.serviceUseCase.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Service not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    service,
	})
}

func (h *ServiceHandler) GetByVehicleID(c *fiber.Ctx) error {
	vehicleID, err := uuid.Parse(c.Params("vehicleId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vehicle ID",
		})
	}

	services, err := h.serviceUseCase.GetByVehicleID(vehicleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    services,
	})
}

func (h *ServiceHandler) GetAll(c *fiber.Ctx) error {
	// Get pagination params from query string
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Validate params
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 5000 {
		limit = 10
	}

	services, total, err := h.serviceUseCase.GetAllWithPagination(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Calculate pagination metadata
	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(fiber.Map{
		"success": true,
		"data":    services,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *ServiceHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid service ID",
		})
	}

	var req dto.UpdateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	service, err := h.serviceUseCase.Update(id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Service updated successfully",
		"data":    service,
	})
}

func (h *ServiceHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid service ID",
		})
	}

	if err := h.serviceUseCase.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Service deleted successfully",
	})
}

func (h *ServiceHandler) CreateQuickTicket(c *fiber.Ctx) error {
	var req dto.QuickTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	service, err := h.serviceUseCase.CreateQuickTicket(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Quick ticket created successfully",
		"data":    service,
	})
}
