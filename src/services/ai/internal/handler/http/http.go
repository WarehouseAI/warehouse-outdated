package http

import (
	"context"
	pg "warehouse/src/internal/database/postgresdb"

	mv "warehouse/src/internal/middleware"
	"warehouse/src/internal/utils/httputils"
	aiCreate "warehouse/src/services/ai/internal/service/ai/create"
	commandCreate "warehouse/src/services/ai/internal/service/command/create"
	"warehouse/src/services/ai/internal/service/command/execute"
	"warehouse/src/services/ai/internal/service/command/get"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type AIServiceProvider struct {
	aiDatabase        *pg.PostgresDatabase[pg.AI]
	commandDatabase   *pg.PostgresDatabase[pg.Command]
	userMiddleware    mv.Middleware
	sessionMiddleware mv.Middleware
	ctx               context.Context
	logger            *logrus.Logger
}

func NewAiAPI(sessionMiddleware mv.Middleware, userMiddleware mv.Middleware, aiDatabase *pg.PostgresDatabase[pg.AI], commandDatabase *pg.PostgresDatabase[pg.Command], logger *logrus.Logger) *AIServiceProvider {
	ctx := context.Background()

	return &AIServiceProvider{
		aiDatabase:        aiDatabase,
		commandDatabase:   commandDatabase,
		userMiddleware:    userMiddleware,
		sessionMiddleware: sessionMiddleware,
		ctx:               ctx,
		logger:            logger,
	}
}

func (pvd *AIServiceProvider) CreateWithoutKeyHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var aiInfo aiCreate.RequestWithoutKey

	if err := c.BodyParser(&aiInfo); err != nil {
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body."))
	}

	ai, err := aiCreate.CreateWithGeneratedKey(&aiInfo, userId, pvd.aiDatabase, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(ai)
}

func (pvd *AIServiceProvider) CreateWithKeyHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var aiInfo aiCreate.RequestWithKey

	if err := c.BodyParser(&aiInfo); err != nil {
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body."))
	}

	ai, err := aiCreate.CreateWithOwnKey(&aiInfo, userId, pvd.aiDatabase, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(ai)
}

func (pvd *AIServiceProvider) AddCommandHandler(c *fiber.Ctx) error {
	var commandCreds commandCreate.Request

	if err := c.BodyParser(&commandCreds); err != nil {
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body."))
	}

	if svcErr := commandCreate.CreateCommand(&commandCreds, pvd.commandDatabase, pvd.logger); svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (pvd *AIServiceProvider) ExecuteCommandHandler(c *fiber.Ctx) error {
	AiID := c.Query("ai_id")
	commandName := c.Query("command_name")

	existCommand, svcErr := get.GetCommand(get.Request{AiID: uuid.FromStringOrNil(AiID), Name: commandName}, pvd.aiDatabase, pvd.logger)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	if existCommand.PayloadType == pg.FormData {
		formData, err := c.MultipartForm()

		if err != nil {
			return c.Status(httputils.InternalError).JSON(httputils.NewErrorResponse(httputils.InternalError, err.Error()))
		}

		response, svcErr := execute.ExecuteFormDataCommand(formData, existCommand, pvd.aiDatabase, pvd.logger)

		if svcErr != nil {
			return c.Status(svcErr.ErrorCode).JSON(svcErr)
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	} else {
		var json map[string]interface{} // не трогать мапу

		if err := c.BodyParser(&json); err != nil {
			statusCode := httputils.BadRequest
			return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body."))
		}

		response, svcErr := execute.ExecuteJSONCommand(json, existCommand, pvd.aiDatabase, pvd.logger)

		if svcErr != nil {
			return c.Status(svcErr.ErrorCode).JSON(svcErr)
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	}
}

// INIT
func (pvd *AIServiceProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(httputils.SetupCORS())
	route := app.Group("/ai")

	route.Post("/create/generate", pvd.sessionMiddleware, pvd.CreateWithoutKeyHandler)
	route.Post("/create/exist", pvd.sessionMiddleware, pvd.CreateWithKeyHandler)
	route.Post("/command/create", pvd.sessionMiddleware, pvd.AddCommandHandler)
	route.Post("/command/execute", pvd.sessionMiddleware, pvd.ExecuteCommandHandler)

	return app
}
