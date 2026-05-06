package handlers

import "github.com/gofiber/fiber/v2"

// HealthCheck responde con el estado del servidor
// GET /health
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "Servidor funcionando correctamente",
	})
}
