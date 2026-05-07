package routes

import (
	"github.com/carrycoders/exposicion/internal/handlers"
	"github.com/carrycoders/exposicion/internal/middleware"
	"github.com/carrycoders/exposicion/internal/repository"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Setup registra todas las rutas de la aplicación
// Recibe la instancia de la BD para inyectarla en los handlers
func Setup(app *fiber.App, db *gorm.DB) {
	// Middleware global
	app.Use(middleware.Logger)

	// Ruta de salud del servidor
	app.Get("/health", handlers.HealthCheck)

	// Wiring: DB → Repository → Handler (patrón de inyección de dependencias)
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Grupo de rutas API v1
	api := app.Group("/api/v1")

	// Rutas de usuarios
	users := api.Group("/users")
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Post("/", userHandler.CreateUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
