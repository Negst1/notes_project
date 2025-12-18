package handlers

import (
	"idm/database/models"
	"idm/database/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandlers(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: authService,
	}
}

func (s *AuthHandler) Register(c *fiber.Ctx) error {
	var registerRequest models.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid request data",
			"details": err.Error(),
		})
	}

	err := s.Register(c)
	if err != nil {
		return err
	}

	if registerRequest.Username == "" || registerRequest.Email == "" || registerRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user data",
		})
	}

	if len(registerRequest.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid password, len < 6",
		})
	}

	return nil
}
