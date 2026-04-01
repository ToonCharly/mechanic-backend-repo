package utils

import "github.com/gofiber/fiber/v2"

// ErrorHandler handles errors in a consistent format
func ErrorHandler(c *fiber.Ctx, err error, status int) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

// SuccessResponse returns a success response
func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}
