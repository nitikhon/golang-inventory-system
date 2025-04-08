package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/util"
)

type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login handles user login and returns access and refresh tokens.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := h.service.Login(loginRequest.Username, loginRequest.Password)
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

// Register handles user registration and returns the created user.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var userRequest entity.User
	if err := c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.Register(userRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
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
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Check if the Authorization header is present
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "access token not found"})
	}

	// Remove the "Bearer " prefix from the token
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid access token format"})
	}

	// Validate the Access Token
	userID, err := util.ValidateAccessToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid access token"})
	}

	// Retrieve the user by ID
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
