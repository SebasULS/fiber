package routes

import (
	"github.com/carrycoders/exposicion/internal/handlers"
	"github.com/carrycoders/exposicion/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// Setup registra todas las rutas de la aplicación
func Setup(app *fiber.App) {
	// Middleware global
	app.Use(middleware.Logger)

	// Ruta de salud del servidor
	app.Get("/health", handlers.HealthCheck)

	// Grupo de rutas API v1
	api := app.Group("/api/v1")

	// Rutas de usuarios
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/", handlers.CreateUser)
	users.Delete("/:id", handlers.DeleteUser)
}
