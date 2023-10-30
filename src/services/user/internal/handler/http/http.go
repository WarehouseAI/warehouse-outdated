package http

import (
	"context"
	"sync"
	pg "warehouse/src/internal/database/postgresdb"
	mw "warehouse/src/internal/middleware"
	u "warehouse/src/internal/utils"
	"warehouse/src/internal/utils/httputils"
	"warehouse/src/services/user/internal/service/update"
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
		return c.Status(fiber.StatusBadRequest).JSON(httputils.NewErrorResponse(httputils.BadRequest, "Invalid request body."))
	}

	updatedUser, err := userUpdate.UpdateUser(updatedFields, user.ID.String(), pvd.userDatabase, pvd.logger)

	if err != nil {
		if err.ErrorCode == httputils.BadRequest {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (pvd *UserServiceProvider) UpdateEmailHandler(c *fiber.Ctx) error {
	var request userUpdate.UpdateEmailRequest
	user := c.Locals("user").(*pg.User)

	if err := c.BodyParser(&request); err != nil || request.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(httputils.NewErrorResponse(httputils.BadRequest, "Invalid request body."))
	}

	key, err := u.GenerateKey(64)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(httputils.NewErrorResponse(httputils.InternalError, err.Error()))
	}

	request.Verified = false
	request.VerificationCode = key

	respch := make(chan *httputils.ErrorResponse, 2)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go update.SendUpdateNotification(wg, respch, pvd.logger, request)
	go update.UpdateEmail(wg, respch, request, user.ID.String(), pvd.userDatabase, pvd.logger)

	wg.Wait()
	close(respch)

	for resp := range respch {
		if resp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(resp)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func (pvd *UserServiceProvider) UpdatePasswordHandler(c *fiber.Ctx) error {
	var request userUpdate.UpdatePasswordRequest
	user := c.Locals("user").(*pg.User)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httputils.NewErrorResponse(httputils.BadRequest, "Invalid request body."))
	}

	if err := userUpdate.UpdatePassword(request, user, pvd.userDatabase, pvd.logger); err != nil {
		if err.ErrorCode == httputils.BadRequest {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (pvd *UserServiceProvider) VerifyUserHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*pg.User)
	verificationCode := c.Params("code")

	if verificationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(httputils.NewErrorResponse(httputils.BadRequest, "Empty verification code"))
	}

	if err := verify.VerifyUserEmail(verify.Request{Verified: true, VerificationCode: &verificationCode}, user, pvd.userDatabase, pvd.logger); err != nil {
		if err.ErrorCode == httputils.BadRequest {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusBadRequest).JSON(err)
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
	route.Post("/update/password", pvd.sessionMiddleware, pvd.userMiddleware, pvd.UpdatePasswordHandler)
	route.Get("/verify/:code", pvd.sessionMiddleware, pvd.userMiddleware, pvd.VerifyUserHandler)

	return app
}
