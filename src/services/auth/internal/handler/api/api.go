package api

import (
	"context"
	"errors"
	"time"
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
		return c.Status(fiber.StatusBadRequest).JSON(im.ErrorResponse{Message: im.BadRequestError.Error()})
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

func (api *AuthAPI) LoginHandler(c *fiber.Ctx) error {
	var loginData m.LoginRequest

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(im.ErrorResponse{Message: im.BadRequestError.Error()})
	}

	session, err := api.svc.Login(context.Background(), &loginData)

	if err != nil && errors.Is(err, im.NotFoundError) {
		return c.Status(fiber.StatusNotFound).JSON(im.ErrorResponse{Message: im.NotFoundError.Error()})
	}

	if err != nil && errors.Is(err, im.BadRequestError) {
		return c.Status(fiber.StatusBadRequest).JSON(im.ErrorResponse{Message: im.BadRequestError.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(im.ErrorResponse{Message: im.InternalError.Error()})
	}

	sessionCookie := new(fiber.Cookie)
	sessionCookie.Name = "id"
	sessionCookie.Value = session.ID
	sessionCookie.Expires = time.Now().Add(session.TTL)

	c.Cookie(sessionCookie)

	return c.SendStatus(fiber.StatusOK)
}

// INIT
func (api *AuthAPI) Init() *fiber.App {
	app := fiber.New()
	route := app.Group("/user")

	route.Post("/create", api.RegisterHandler)
	route.Post("/login", api.LoginHandler)

	return app
}
