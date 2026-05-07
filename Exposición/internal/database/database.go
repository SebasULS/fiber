package database

import (
	"log"

	"github.com/carrycoders/exposicion/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia global de la base de datos
var DB *gorm.DB

// Connect abre la conexión con SQLite y ejecuta las migraciones automáticas
func Connect(dsn string) {
	var err error

	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// AutoMigrate crea o actualiza la tabla 'users' según el modelo
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Error en la migración automática: %v", err)
	}

	log.Println("Base de datos SQLite conectada y migrada correctamente")
}
