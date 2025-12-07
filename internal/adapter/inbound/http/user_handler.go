package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/core/service/validation"
	"github.com/nitikhon/golang-inventory-system/internal/util"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate and normalize user input
	if err := validation.ValidateAndNormalizeUser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdUser, err := h.service.CreateUser(&user)
	if err != nil {
		if err.Error() == "a user with the provided credentials already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// Update modifies an existing user.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	var user entity.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user ID is required for update"})
	}

	// Validate and normalize user input
	if err := validation.ValidateAndNormalizeUser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedUser, err := h.service.UpdateUser(&user)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
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

	err = h.service.DeleteUser(uint(id))
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

	user, err := h.service.GetUserByID(uint(id))
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

	// Validate and normalize username
	username, err := validation.ValidateAndNormalizeUsername(username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
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

	// Validate and normalize email
	email, err := validation.ValidateAndNormalizeEmail(email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
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

	// Validate and normalize phone
	phone, err := validation.ValidateAndNormalizePhone(phone)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.GetUserByPhone(phone)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// Login handles user login and returns access and refresh tokens.
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// Validate and normalize username
	username, err := validation.ValidateAndNormalizeUsername(loginRequest.Username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := h.service.Login(username, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// Set the Refresh Token in an HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   string(util.GetEnvOrPanic("COOKIE_SECURE")) == "true", // Set to true in production
		SameSite: "Strict",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

// RefreshToken handles the HTTP request to refresh authentication tokens.
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	// Retrieve the Refresh Token from the HTTP-only cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "refresh token not found"})
	}

	// Validate the Refresh Token and generate new tokens
	tokens, err := h.service.RefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	// Set the new Refresh Token in the HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Secure:   string(util.GetEnvOrPanic("COOKIE_SECURE")) == "true", // Set to true in production
		SameSite: "Strict",
		Path:     "/",
	})

	// Return the new Access Token in the response body
	return c.JSON(fiber.Map{
		"access_token": tokens.AccessToken,
	})
}

// Me retrieves the current user's information using the Access Token.
// It validates the token and returns the user data if valid.
func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Retrieve the user by ID
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// May exclude some field from the response later
	return c.JSON(user)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Get user ID from context (set by AuthMiddleware)
	userID := c.Locals("user_id").(uint)

	// Clear refresh token in DB
	err := h.service.Logout(userID)
	if err != nil {
		if err.Error() == "user already logged out" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Already logged out",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Clear the refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
	})

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}
