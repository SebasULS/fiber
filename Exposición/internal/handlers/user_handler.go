package handlers

import (
	"strconv"

	"github.com/carrycoders/exposicion/internal/models"
	"github.com/gofiber/fiber/v2"
)

// Simulación de base de datos en memoria
var users = []models.User{
	{ID: 1, Name: "Alice García", Email: "alice@example.com", Age: 25},
	{ID: 2, Name: "Bob Martínez", Email: "bob@example.com", Age: 30},
}

// GetUsers devuelve todos los usuarios
// GET /api/users
func GetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data":  users,
		"total": len(users),
	})
}

// GetUser devuelve un usuario por ID
// GET /api/users/:id
func GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	for _, u := range users {
		if u.ID == id {
			return c.JSON(u)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Usuario no encontrado",
	})
}

// CreateUser crea un nuevo usuario
// POST /api/users
func CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cuerpo de petición inválido",
		})
	}

	newUser := models.User{
		ID:    len(users) + 1,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	users = append(users, newUser)

	return c.Status(fiber.StatusCreated).JSON(newUser)
}

// DeleteUser elimina un usuario por ID
// DELETE /api/users/:id
func DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			return c.Status(fiber.StatusNoContent).Send(nil)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Usuario no encontrado",
	})
}
