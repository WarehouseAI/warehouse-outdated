package http

import (
	"errors"
	"warehouse/gen"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/dto"
	mv "warehouse/src/internal/middleware"
	u "warehouse/src/internal/utils"
	svc "warehouse/src/services/auth/internal/service"
	m "warehouse/src/services/auth/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ApiProvider struct {
	sessionDatabase   *r.RedisDatabase
	sessionMiddleware *mv.SessionMiddlewareProvider
	authService       *svc.AuthServiceProvider
}

func NewAuthAPI(sessionDatabase *r.RedisDatabase, logger *logrus.Logger) *ApiProvider {
	authService := svc.NewAuthService(sessionDatabase, logger)
	sessionMiddleware := mv.NewSessionMiddleware(sessionDatabase, logger)

	return &ApiProvider{
		sessionDatabase:   sessionDatabase,
		sessionMiddleware: sessionMiddleware,
		authService:       authService,
	}
}

func (api *ApiProvider) RegisterHandler(c *fiber.Ctx) error {
	var userInfo gen.CreateUserMsg

	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	userId, err := api.authService.Register(&userInfo)

	if err != nil && errors.Is(err, dto.ExistError) {
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(userId)
}

func (api *ApiProvider) LoginHandler(c *fiber.Ctx) error {
	var loginData m.LoginRequest

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	session, err := api.authService.Login(&loginData)

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
		Name:  "sessionId",
		Value: session.ID,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (api *ApiProvider) LogoutHandler(c *fiber.Ctx) error {
	sessionId := c.Cookies("sessionId")

	if err := api.authService.Logout(sessionId); err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (api *ApiProvider) WhoAmIHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

// INIT
func (api *ApiProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(u.SetupCORS())
	route := app.Group("/auth")

	route.Post("/register", api.RegisterHandler)
	route.Post("/login", api.LoginHandler)
	route.Delete("/logout", api.sessionMiddleware.Session, api.LogoutHandler)
	route.Get("/whoami", api.sessionMiddleware.Session, api.WhoAmIHandler)

	return app
}
