package commands

import (
	"warehouseai/ai/adapter/grpc/client/auth"
	"warehouseai/ai/dataservice/aidata"
	"warehouseai/ai/dataservice/commanddata"
	e "warehouseai/ai/errors"
	"warehouseai/ai/service/command"
	"warehouseai/ai/service/command/execute"

	"github.com/gofiber/fiber/v2"
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
	aiID := c.Query("ai_id")
	commandName := c.Query("command_name")

	request := execute.ExecuteCommandRequest{
		AiID:        aiID,
		CommandName: commandName,
		Raw:         c.Body(),
		ContentType: c.Get("Content-Type"),
	}

	response, err := execute.ExecuteCommand(request, h.AiDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	for key, value := range response.Headers {
		c.Response().Header.Add(key, value)
	}

	return c.Status(response.Status).Send(response.Raw.Bytes())
}
