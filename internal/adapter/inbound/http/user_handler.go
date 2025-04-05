package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Create adds a new user.
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var user entity.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdUser, err := h.service.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(createdUser)
}

// Update modifies an existing user.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	var user entity.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedUser, err := h.service.Update(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedUser)
}

// Delete removes a user by their ID.
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.service.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetAllUsers retrieves all users.
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// GetUserByID retrieves a user by their ID.
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// GetUserByUsername retrieves a user by their username.
func (h *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username is required"})
	}

	user, err := h.service.GetUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// GetUserByEmail retrieves a user by their email.
func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// GetUserByPhone retrieves a user by their phone number.
func (h *UserHandler) GetUserByPhone(c *fiber.Ctx) error {
	phone := c.Params("phone")
	if phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "phone is required"})
	}

	user, err := h.service.GetUserByPhone(phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}