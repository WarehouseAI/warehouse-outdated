package http

import (
	"context"
	"sync"
	"warehouse/gen"
	pg "warehouse/src/internal/database/postgresdb"
	mw "warehouse/src/internal/middleware"
	u "warehouse/src/internal/utils"
	"warehouse/src/internal/utils/httputils"
	gw "warehouse/src/services/user/internal/gateway"
	"warehouse/src/services/user/internal/service/get"
	"warehouse/src/services/user/internal/service/update"
	userUpdate "warehouse/src/services/user/internal/service/update"
	"warehouse/src/services/user/internal/service/verify"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserServiceProvider struct {
	userDatabase      *pg.PostgresDatabase[pg.User]
	userMiddleware    mw.Middleware
	aiGateway         *gw.AiGrpcConnection
	sessionMiddleware mw.Middleware
	logger            *logrus.Logger
	ctx               context.Context
}

func NewUserPublicAPI(userDatabase *pg.PostgresDatabase[pg.User], userMiddleware mw.Middleware, sessionMiddleware mw.Middleware, logger *logrus.Logger) *UserServiceProvider {
	ctx := context.Background()
	ai := gw.NewAiGrpcConnection("ai-service:8021")

	return &UserServiceProvider{
		userDatabase:      userDatabase,
		sessionMiddleware: sessionMiddleware,
		aiGateway:         ai,
		userMiddleware:    userMiddleware,
		logger:            logger,
		ctx:               ctx,
	}
}

func (pvd *UserServiceProvider) UpdateHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var updatedFields userUpdate.UpdateUserRequest

	if err := c.BodyParser(&updatedFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httputils.NewErrorResponse(httputils.BadRequest, "Invalid request body."))
	}

	updatedUser, err := userUpdate.UpdateUser(updatedFields, "id", userId, pvd.userDatabase, pvd.logger)

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
	userId := c.Locals("userId").(string)

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
	go update.UpdateUserEmail(wg, respch, request, userId, pvd.userDatabase, pvd.logger)

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

	if err := userUpdate.UpdateUserPassword(request, user, pvd.userDatabase, pvd.logger); err != nil {
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

func (pvd *UserServiceProvider) AddFavoriteAIHandler(c *fiber.Ctx) error {
	var request *gen.GetAiByIdMsg
	user := c.Locals("user").(*pg.User)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httputils.NewErrorResponse(httputils.BadRequest, "Invalid request body."))
	}

	if svcErr := update.AddUserFavoriteAi(request, user, pvd.userDatabase, pvd.aiGateway, pvd.logger); svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (pvd *UserServiceProvider) RemoveFavoriteAIHandler(c *fiber.Ctx) error {
	var request *gen.GetAiByIdMsg
	user := c.Locals("user").(*pg.User)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(httputils.NewErrorResponse(httputils.BadRequest, "Invalid request body."))
	}

	if svcErr := update.RemoveUserFavoriteAi(request, user, pvd.userDatabase, pvd.aiGateway, pvd.logger); svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (pvd *UserServiceProvider) GetFavoriteAIHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	favoriteAi, svcErr := get.GetUserFavoriteAi(userId, pvd.userDatabase, pvd.logger)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.Status(fiber.StatusOK).JSON(favoriteAi)
}

func (pvd *UserServiceProvider) GetUserById(c *fiber.Ctx) error {
	userId := c.Query("id")

	user, err := get.GetById(&gen.GetUserByIdMsg{Id: userId}, pvd.userDatabase, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// INIT
func (pvd *UserServiceProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(httputils.SetupCORS())
	route := app.Group("/user")

	route.Patch("/update", pvd.sessionMiddleware, pvd.UpdateHandler)
	route.Patch("/update/email", pvd.sessionMiddleware, pvd.UpdateEmailHandler)
	route.Patch("/update/password", pvd.sessionMiddleware, pvd.userMiddleware, pvd.UpdatePasswordHandler)
	route.Patch("/favorites/add", pvd.sessionMiddleware, pvd.userMiddleware, pvd.AddFavoriteAIHandler)
	route.Patch("/favorites/remove", pvd.sessionMiddleware, pvd.userMiddleware, pvd.RemoveFavoriteAIHandler)
	route.Get("/get", pvd.GetUserById)
	route.Get("/favorites", pvd.sessionMiddleware, pvd.GetFavoriteAIHandler)
	route.Get("/verify/:code", pvd.sessionMiddleware, pvd.userMiddleware, pvd.VerifyUserHandler)

	return app
}
