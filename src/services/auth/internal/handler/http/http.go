package http

import (
	"context"
	"errors"
	"warehouse/gen"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/dto"
	mv "warehouse/src/internal/middleware"
	"warehouse/src/internal/utils/httputils"
	gw "warehouse/src/services/auth/internal/gateway"
	"warehouse/src/services/auth/internal/service/login"
	"warehouse/src/services/auth/internal/service/logout"
	"warehouse/src/services/auth/internal/service/register"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthServiceProvider struct {
	sessionDatabase   *r.RedisDatabase
	sessionMiddleware mv.Middleware
	userGateway       *gw.UserGrpcConnection
	logger            *logrus.Logger
	ctx               context.Context
}

func NewAuthAPI(sessionDatabase *r.RedisDatabase, sessionMiddleware mv.Middleware, logger *logrus.Logger) *AuthServiceProvider {
	ctx := context.Background()
	userGateway := gw.NewUserGrpcConnection("user-service:8001")

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
	var creds login.Request

	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	session, err := login.Login(&creds, pvd.userGateway, pvd.sessionDatabase, pvd.logger, pvd.ctx)

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
		Name:     "sessionId",
		Value:    session.ID,
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
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
func (pvd *AuthServiceProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(httputils.SetupCORS())
	route := app.Group("/auth")

	route.Post("/register", pvd.RegisterHandler)
	route.Post("/login", pvd.LoginHandler)
	route.Delete("/logout", pvd.sessionMiddleware, pvd.LogoutHandler)
	route.Get("/whoami", pvd.sessionMiddleware, pvd.WhoAmIHandler)

	return app
}
