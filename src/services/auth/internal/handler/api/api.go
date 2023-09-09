package api

import (
	"context"
	"errors"
	"warehouse/gen"
	im "warehouse/src/internal/models"
	m "warehouse/src/services/auth/pkg/model"

	"github.com/gofiber/fiber/v2"
)

type AuthAPI struct {
	svc m.AuthService
}

func NewAuthAPI(svc m.AuthService) *AuthAPI {
	return &AuthAPI{
		svc: svc,
	}
}

func (api *AuthAPI) RegisterHandler(c *fiber.Ctx) error {
	var userInfo gen.CreateUserRequest

	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(im.ErrorResponse{Message: "Invalid request body"})
	}

	userId, err := api.svc.Register(context.Background(), &userInfo)

	if err != nil && errors.Is(err, im.InternalError) {
		return c.Status(fiber.StatusInternalServerError).JSON(im.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(im.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(userId)
}

// INIT
func (api *AuthAPI) Init() *fiber.App {
	app := fiber.New()
	route := app.Group("/user")

	route.Post("/create", api.RegisterHandler)

	return app
}
