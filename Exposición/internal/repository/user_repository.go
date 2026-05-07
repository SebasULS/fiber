package repository

import (
	"github.com/carrycoders/exposicion/internal/models"
	"gorm.io/gorm"
)

// UserRepository encapsula todas las operaciones de la tabla users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository crea un nuevo repositorio inyectando la conexión de BD
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll obtiene todos los usuarios de la BD
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	return users, result.Error
}

// FindByID obtiene un usuario por su ID, devuelve error si no existe
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	return &user, result.Error
}

// Create inserta un nuevo usuario en la BD
func (r *UserRepository) Create(req *models.CreateUserRequest) (*models.User, error) {
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	result := r.db.Create(&user)
	return &user, result.Error
}

// Update modifica un usuario existente por ID
func (r *UserRepository) Update(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if result := r.db.First(&user, id); result.Error != nil {
		return nil, result.Error
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age
	result := r.db.Save(&user)
	return &user, result.Error
}

// Delete elimina un usuario por ID
func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}
