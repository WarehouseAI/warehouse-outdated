package http

import (
	"context"
	"errors"
	dbm "warehouse/src/internal/db/models"
	"warehouse/src/internal/dto"
	mv "warehouse/src/internal/middleware"
	svc "warehouse/src/services/ai/internal/service"
	m "warehouse/src/services/ai/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type APIInstance struct {
	svc svc.AIService
	sMw *mv.SessionMiddleware
	uMw *mv.UserMiddleware
}

func NewAiAPI(service svc.AIService, sessionMiddleware *mv.SessionMiddleware, userMiddleware *mv.UserMiddleware) *APIInstance {
	return &APIInstance{
		svc: service,
		sMw: sessionMiddleware,
		uMw: userMiddleware,
	}
}

func (api *APIInstance) CreateHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*dbm.User)
	var aiInfo m.CreateAIRequest

	if err := c.BodyParser(&aiInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	apiInfo, err := api.svc.CreateWithGeneratedKey(context.Background(), &aiInfo, user)

	if err != nil && errors.Is(err, dto.InternalError) {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(apiInfo)
}

func (api *APIInstance) CreateWithKey(c *fiber.Ctx) error {
	user := c.Locals("user").(*dbm.User)
	var aiInfo m.CreateAIRequest

	if err := c.BodyParser(&aiInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	apiInfo, err := api.svc.CreateWithOwnKey(context.Background(), &aiInfo, user)

	if err != nil && errors.Is(err, dto.InternalError) {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(apiInfo)
}

func (api *APIInstance) AddCommandHandler(c *fiber.Ctx) error {
	var commandInfo m.AddCommandRequest

	if err := c.BodyParser(&commandInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	existCommand, err := api.svc.GetCommand(context.Background(), commandInfo.AiID.String(), commandInfo.Name)

	if existCommand != nil {
		statusCode := fiber.StatusConflict
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.ExistError.Error()})
	}

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	if err := api.svc.AddCommand(context.Background(), &commandInfo); err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (api *APIInstance) ExecuteCommandHandler(c *fiber.Ctx) error {
	aiID := c.Query("ai_id")
	commandName := c.Query("command_name")

	existCommand, err := api.svc.GetCommand(context.Background(), aiID, commandName)

	if existCommand == nil {
		statusCode := fiber.StatusNotFound
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	if existCommand.PayloadType == dbm.FormData {
		formData, err := c.MultipartForm()

		if err != nil {
			statusCode := fiber.StatusBadRequest
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		response, err := api.svc.ExecuteFormDataCommand(context.Background(), formData, existCommand)

		if err != nil {
			statusCode := fiber.StatusInternalServerError
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	} else {
		var json map[string]interface{} // оставить мап т.к. дальше отправляю в нейронку его

		if err := c.BodyParser(&json); err != nil {
			statusCode := fiber.StatusBadRequest
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		response, err := api.svc.ExecuteJSONCommand(context.Background(), json, existCommand)

		if err != nil {
			statusCode := fiber.StatusInternalServerError
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	}
}

// INIT
func (api *APIInstance) Init() *fiber.App {
	app := fiber.New()
	route := app.Group("/ai")

	route.Post("/create", api.sMw.Session, api.uMw.User, api.CreateHandler) // Combine to one
	route.Post("/createWith", api.sMw.Session, api.uMw.User, api.CreateWithKey)
	route.Post("/command/create", api.sMw.Session, api.AddCommandHandler)
	route.Post("/command/execute", api.sMw.Session, api.ExecuteCommandHandler)

	return app
}
