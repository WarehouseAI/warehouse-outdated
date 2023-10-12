package http

import (
	"context"
	"errors"
	"fmt"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"
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
	user := c.Locals("user").(*pg.User)
	var aiInfo aiCreate.RequestWithoutKey

	if err := c.BodyParser(&aiInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	ai, err := aiCreate.CreateWithGeneratedKey(&aiInfo, user, pvd.aiDatabase, pvd.logger, pvd.ctx)

	if err != nil && errors.Is(err, dto.InternalError) {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(ai)
}

func (pvd *AIServiceProvider) CreateWithKeyHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*pg.User)
	var aiInfo aiCreate.RequestWithKey

	if err := c.BodyParser(&aiInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	ai, err := aiCreate.CreateWithOwnKey(&aiInfo, user, pvd.aiDatabase, pvd.logger, pvd.ctx)

	if err != nil && errors.Is(err, dto.InternalError) {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(ai)
}

func (pvd *AIServiceProvider) AddCommandHandler(c *fiber.Ctx) error {
	var commandCreds commandCreate.Request

	if err := c.BodyParser(&commandCreds); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	if err := commandCreate.CreateCommand(&commandCreds, pvd.commandDatabase, pvd.logger); err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (pvd *AIServiceProvider) ExecuteCommandHandler(c *fiber.Ctx) error {
	AiID := c.Query("ai_id")
	commandName := c.Query("command_name")

	existCommand, err := get.GetCommand(get.Request{AiID: uuid.FromStringOrNil(AiID), Name: commandName}, pvd.commandDatabase, pvd.logger)

	if existCommand == nil {
		statusCode := fiber.StatusNotFound
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	if existCommand.PayloadType == pg.FormData {
		formData, err := c.MultipartForm()

		if err != nil {
			statusCode := fiber.StatusBadRequest
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		response, err := execute.ExecuteFormDataCommand(formData, existCommand, pvd.aiDatabase, pvd.logger)

		if err != nil {
			statusCode := fiber.StatusInternalServerError
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	} else {
		var json map[string]interface{} // не трогать мапу

		if err := c.BodyParser(&json); err != nil {
			statusCode := fiber.StatusBadRequest
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		response, err := execute.ExecuteJSONCommand(json, existCommand, pvd.aiDatabase, pvd.logger)

		if err != nil {
			statusCode := fiber.StatusInternalServerError
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	}
}

// INIT
func (pvd *AIServiceProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(httputils.SetupCORS())
	route := app.Group("/ai")

	route.Post("/create/generate", pvd.sessionMiddleware, pvd.userMiddleware, pvd.CreateWithoutKeyHandler)
	route.Post("/create/exist", pvd.sessionMiddleware, pvd.userMiddleware, pvd.CreateWithKeyHandler)
	route.Post("/command/create", pvd.sessionMiddleware, pvd.AddCommandHandler)
	route.Post("/command/execute", pvd.sessionMiddleware, pvd.ExecuteCommandHandler)

	return app
}
