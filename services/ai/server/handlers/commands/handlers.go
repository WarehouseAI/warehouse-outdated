package commands

import (
	"warehouseai/ai/adapter/grpc/client/auth"
	"warehouseai/ai/dataservice/aidata"
	"warehouseai/ai/dataservice/commanddata"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
	"warehouseai/ai/service/command"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	CommandDB  *commanddata.Database
	AiDB       *aidata.Database
	Logger     *logrus.Logger
	AuthClient *auth.AuthGrpcClient
}

func (h *Handler) CreateCommandHandler(c *fiber.Ctx) error {
	var commandCreds command.CreateCommandRequest

	if err := c.BodyParser(&commandCreds); err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, "Invalid request body.")
		return c.Status(response.ErrorCode).JSON(response)
	}

	if svcErr := command.CreateCommand(&commandCreds, h.CommandDB, h.Logger); svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) ExecuteCommandHandler(c *fiber.Ctx) error {
	AiID := c.Query("ai_id")
	commandName := c.Query("command_name")

	getCommandRequest := command.GetCommandRequest{
		AiID: uuid.FromStringOrNil(AiID),
		Name: commandName,
	}

	existCommand, svcErr := command.GetCommand(getCommandRequest, h.AiDB, h.Logger)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	if existCommand.Payload.PayloadType == m.FormData {
		formData, err := c.MultipartForm()

		if err != nil {
			response := e.NewErrorResponse(e.HttpInternalError, err.Error())
			return c.Status(response.ErrorCode).JSON(response)
		}

		response, svcErr := command.ExecuteFormDataCommand(formData, existCommand, h.AiDB, h.Logger)

		if svcErr != nil {
			return c.Status(svcErr.ErrorCode).JSON(svcErr)
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	} else {
		var json map[string]interface{} // не трогать мапу

		if err := c.BodyParser(&json); err != nil {
			response := e.NewErrorResponse(e.HttpInternalError, "Invalid request body.")
			return c.Status(response.ErrorCode).JSON(response)
		}

		response, svcErr := command.ExecuteJSONCommand(json, existCommand, h.AiDB, h.Logger)

		if svcErr != nil {
			return c.Status(svcErr.ErrorCode).JSON(svcErr)
		}

		return c.Status(fiber.StatusOK).Send(response.Bytes())
	}
}
