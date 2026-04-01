package handlers

import (
	"mechanic-backend/internal/application/dto"
	"mechanic-backend/internal/application/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VehicleHandler struct {
	vehicleUseCase *usecases.VehicleUseCase
}

func NewVehicleHandler(vehicleUseCase *usecases.VehicleUseCase) *VehicleHandler {
	return &VehicleHandler{vehicleUseCase: vehicleUseCase}
}

func (h *VehicleHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateVehicleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	vehicle, err := h.vehicleUseCase.Create(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Vehicle created successfully",
		"data":    vehicle,
	})
}

func (h *VehicleHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vehicle ID",
		})
	}

	vehicle, err := h.vehicleUseCase.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Vehicle not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    vehicle,
	})
}

func (h *VehicleHandler) GetAll(c *fiber.Ctx) error {
	vehicles, err := h.vehicleUseCase.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    vehicles,
	})
}

func (h *VehicleHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vehicle ID",
		})
	}

	var req dto.UpdateVehicleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	vehicle, err := h.vehicleUseCase.Update(id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Vehicle updated successfully",
		"data":    vehicle,
	})
}

func (h *VehicleHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid vehicle ID",
		})
	}

	if err := h.vehicleUseCase.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Vehicle deleted successfully",
	})
}
