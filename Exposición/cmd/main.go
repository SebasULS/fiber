package main

import (
	"log"

	"github.com/carrycoders/exposicion/internal/config"
	"github.com/carrycoders/exposicion/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Cargar configuración de la aplicación
	cfg := config.Load()

	// Crear instancia de Fiber
	app := fiber.New(fiber.Config{
		AppName: "CarryCoders API v1.0",
	})

	// Registrar todas las rutas
	routes.Setup(app)

	// Iniciar servidor
	log.Printf("Servidor corriendo en http://localhost%s\n", cfg.Port)
	log.Fatal(app.Listen(cfg.Port))
}
