package http

import (
	"context"
	"errors"
	"warehouse/gen"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/dto"
	mv "warehouse/src/internal/middleware"
	u "warehouse/src/internal/utils"
	gw "warehouse/src/services/auth/internal/gateway"
	"warehouse/src/services/auth/internal/service/login"
	"warehouse/src/services/auth/internal/service/logout"
	"warehouse/src/services/auth/internal/service/register"
	m "warehouse/src/services/auth/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthServiceProvider struct {
	sessionDatabase   *r.RedisDatabase
	sessionMiddleware *mv.SessionMiddlewareProvider
	userGateway       *gw.UserGrpcConnection
	logger            *logrus.Logger
	ctx               context.Context
}

func NewAuthAPI(sessionDatabase *r.RedisDatabase, logger *logrus.Logger) *AuthServiceProvider {
	ctx := context.Background()
	userGateway := gw.NewUserGrpcConnection("user-service:8001")
	sessionMiddleware := mv.NewSessionMiddleware(sessionDatabase, logger)

	return &AuthServiceProvider{
		userGateway:       userGateway,
		sessionDatabase:   sessionDatabase,
		sessionMiddleware: sessionMiddleware,
		ctx:               ctx,
		logger:            logger,
	}
}

func (pvd *AuthServiceProvider) RegisterHandler(c *fiber.Ctx) error {
	var userInfo gen.CreateUserMsg

	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	userId, err := register.Register(&userInfo, pvd.userGateway, pvd.logger, pvd.ctx)

	if err != nil && errors.Is(err, dto.ExistError) {
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(userId)
}

func (pvd *AuthServiceProvider) LoginHandler(c *fiber.Ctx) error {
	var loginData m.LoginRequest

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	session, err := login.Login(&loginData, pvd.userGateway, pvd.sessionDatabase, pvd.logger, pvd.ctx)

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

func (pvd *AuthServiceProvider) LogoutHandler(c *fiber.Ctx) error {
	sessionId := c.Cookies("sessionId")

	if err := logout.Logout(sessionId, pvd.sessionDatabase, pvd.logger, pvd.ctx); err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (api *AuthServiceProvider) WhoAmIHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

// INIT
func (api *AuthServiceProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(u.SetupCORS())
	route := app.Group("/auth")

	route.Post("/register", api.RegisterHandler)
	route.Post("/login", api.LoginHandler)
	route.Delete("/logout", api.sessionMiddleware.Session, api.LogoutHandler)
	route.Get("/whoami", api.sessionMiddleware.Session, api.WhoAmIHandler)

	return app
}
