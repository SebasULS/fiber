package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger es un middleware personalizado que registra cada petición HTTP
func Logger(c *fiber.Ctx) error {
	start := time.Now()

	// Continuar con el siguiente handler
	err := c.Next()

	// Registrar después de que la respuesta fue enviada
	log.Printf(
		"[%s] %s %s - %d (%s)",
		time.Now().Format("2006-01-02 15:04:05"),
		c.Method(),
		c.Path(),
		c.Response().StatusCode(),
		time.Since(start),
	)

	return err
}
