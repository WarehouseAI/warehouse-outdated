package http

import (
	"context"
	"warehouse/gen"
	r "warehouse/src/internal/database/redisdb"
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
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body"))
	}

	userId, err := register.Register(&userInfo, pvd.userGateway, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(userId)
}

func (pvd *AuthServiceProvider) LoginHandler(c *fiber.Ctx) error {
	var creds login.Request

	if err := c.BodyParser(&creds); err != nil {
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body"))
	}

	session, err := login.Login(&creds, pvd.userGateway, pvd.sessionDatabase, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
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

	if sessionId == "" {
		return c.Status(httputils.Unauthorized).JSON(httputils.NewErrorResponse(httputils.Unauthorized, "Empty session key."))
	}

	if err := logout.Logout(sessionId, pvd.sessionDatabase, pvd.logger, pvd.ctx); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}
	c.ClearCookie("sessionId")

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
	route.Delete("/logout", pvd.LogoutHandler)
	route.Get("/whoami", pvd.sessionMiddleware, pvd.WhoAmIHandler)

	return app
}
