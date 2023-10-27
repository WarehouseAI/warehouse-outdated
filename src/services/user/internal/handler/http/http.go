package http

import (
	"context"
	"errors"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"
	mw "warehouse/src/internal/middleware"
	"warehouse/src/internal/utils/httputils"
	userUpdate "warehouse/src/services/user/internal/service/update"
	"warehouse/src/services/user/internal/service/verify"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserServiceProvider struct {
	userDatabase      *pg.PostgresDatabase[pg.User]
	userMiddleware    mw.Middleware
	sessionMiddleware mw.Middleware
	logger            *logrus.Logger
	ctx               context.Context
}

func NewUserPublicAPI(userDatabase *pg.PostgresDatabase[pg.User], userMiddleware mw.Middleware, sessionMiddleware mw.Middleware, logger *logrus.Logger) *UserServiceProvider {
	ctx := context.Background()

	return &UserServiceProvider{
		userDatabase:      userDatabase,
		sessionMiddleware: sessionMiddleware,
		userMiddleware:    userMiddleware,
		logger:            logger,
		ctx:               ctx,
	}
}

func (pvd *UserServiceProvider) UpdateHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*pg.User)
	var updatedFields userUpdate.UpdateUserRequest

	if err := c.BodyParser(&updatedFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	updatedUser, err := userUpdate.UpdateUser(updatedFields, user.ID.String(), pvd.userDatabase, pvd.logger)

	if err != nil && errors.Is(err, dto.BadRequestError) {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (pvd *UserServiceProvider) UpdateEmailHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*pg.User)
	var newEmail userUpdate.UpdateEmailRequest

	if err := c.BodyParser(&newEmail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	if err := userUpdate.UpdateEmail(newEmail, user.ID.String(), pvd.userDatabase, pvd.logger); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (pvd *UserServiceProvider) VerifyUserHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*pg.User)
	verificationCode := c.Params("code")

	if verificationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	if err := verify.VerifyUserEmail(verify.Request{Verified: true, VerificationCode: verificationCode}, user, pvd.userDatabase, pvd.logger); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// INIT
func (pvd *UserServiceProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(httputils.SetupCORS())
	route := app.Group("/user")

	route.Post("/update", pvd.sessionMiddleware, pvd.userMiddleware, pvd.UpdateHandler)
	route.Post("/update/email", pvd.sessionMiddleware, pvd.userMiddleware, pvd.UpdateEmailHandler)
	route.Get("/verify/:code", pvd.sessionMiddleware, pvd.userMiddleware, pvd.VerifyUserHandler)

	return app
}
