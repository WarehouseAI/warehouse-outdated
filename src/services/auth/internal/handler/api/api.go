package api

import (
	"context"
	"errors"
	"warehouse/gen"
	"warehouse/src/internal/dto"
	mv "warehouse/src/internal/middleware"
	svc "warehouse/src/services/auth/internal/service/auth"
	m "warehouse/src/services/auth/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type APIInstance struct {
	svc svc.AuthService
	mvs mv.MiddlewareService
}

func NewAuthAPI(svc svc.AuthService, mvs mv.MiddlewareService) *APIInstance {
	return &APIInstance{
		svc: svc,
		mvs: mvs,
	}
}

func (api *APIInstance) RegisterHandler(c *fiber.Ctx) error {
	var userInfo gen.CreateUserRequest

	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	userId, err := api.svc.Register(context.Background(), &userInfo)

	if err != nil && errors.Is(err, dto.InternalError) {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(userId)
}

func (api *APIInstance) LoginHandler(c *fiber.Ctx) error {
	var loginData m.LoginRequest

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	session, err := api.svc.Login(context.Background(), &loginData)

	if err != nil && errors.Is(err, dto.NotFoundError) {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Message: dto.NotFoundError.Error()})
	}

	if err != nil && errors.Is(err, dto.BadRequestError) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: dto.InternalError.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "id",
		Value: session.ID,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (api *APIInstance) LogoutHandler(c *fiber.Ctx) error {
	sessionid := c.Cookies("id")

	if err := api.svc.Logout(context.Background(), sessionid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: dto.InternalError.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// INIT
func (api *APIInstance) Init() *fiber.App {
	app := fiber.New()
	route := app.Group("/auth")

	route.Post("/register", api.RegisterHandler)
	route.Post("/login", api.LoginHandler)
	route.Delete("/logout", api.mvs.Session, api.LogoutHandler)

	return app
}
