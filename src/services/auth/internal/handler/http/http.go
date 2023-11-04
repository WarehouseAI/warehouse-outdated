package http

import (
	"context"
	"warehouse/gen"
	r "warehouse/src/internal/database/redisdb"
	mv "warehouse/src/internal/middleware"
	"warehouse/src/internal/s3"
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
	s3                *s3.S3Storage
	sessionMiddleware mv.Middleware
	userGateway       *gw.UserGrpcConnection
	logger            *logrus.Logger
	ctx               context.Context
}

func NewAuthAPI(sessionDatabase *r.RedisDatabase, sessionMiddleware mv.Middleware, logger *logrus.Logger, s3 *s3.S3Storage) *AuthServiceProvider {
	ctx := context.Background()
	userGateway := gw.NewUserGrpcConnection("user-service:8001")

	return &AuthServiceProvider{
		userGateway:       userGateway,
		sessionDatabase:   sessionDatabase,
		sessionMiddleware: sessionMiddleware,
		s3:                s3,
		ctx:               ctx,
		logger:            logger,
	}
}

func (pvd *AuthServiceProvider) RegisterHandler(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		response := httputils.NewErrorResponse(httputils.BadRequest, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	username := form.Value["username"][0]
	rawPicture, err := c.FormFile("picture")

	if err != nil {
		response := httputils.NewErrorResponse(httputils.InternalError, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	picture, err := rawPicture.Open()

	if err != nil {
		response := httputils.NewErrorResponse(httputils.InternalError, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	defer picture.Close()

	link, svcErr := register.UploadAvatar(picture, username, pvd.logger, pvd.s3)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	userInfo := &gen.CreateUserMsg{
		Username:  username,
		Firstname: form.Value["firstname"][0],
		Lastname:  form.Value["lastname"][0],
		Password:  form.Value["password"][0],
		Email:     form.Value["email"][0],
		Picture:   link,
	}

	userId, svcErr := register.Register(userInfo, pvd.userGateway, pvd.logger, pvd.ctx)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
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
