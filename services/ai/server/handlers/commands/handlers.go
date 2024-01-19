package commands

import (
	"mime/multipart"
	"warehouseai/ai/adapter/grpc/client/auth"
	"warehouseai/ai/dataservice/psql/aidata"
	"warehouseai/ai/dataservice/psql/commanddata"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
	"warehouseai/ai/service/command/create"
	"warehouseai/ai/service/command/execute"
	"warehouseai/ai/service/command/get"

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
	var commandCreds create.CreateCommandRequest

	if err := c.BodyParser(&commandCreds); err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	if svcErr := create.CreateCommand(&commandCreds, h.CommandDB, h.Logger); svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// Решил логику определения типа запроса, для корректного парсинга, перенести сюда.
// Все таки она не относится к бизнес-логике, а скорее к логике обработки запросов, и код в общем становится чище.
func (h *Handler) ExecuteCommandHandler(c *fiber.Ctx) error {
	aiID := c.Query("ai_id")
	commandName := c.Query("command_name")

	existCommandInfo, err := get.GetCommand(get.GetCommandRequest{AiID: aiID, Name: commandName}, h.AiDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	if existCommandInfo.Command.PayloadType == string(m.FormData) {
		formPayload, err := c.MultipartForm()

		if err != nil {
			resp := e.NewErrorResponse(e.HttpInternalError, err.Error())
			return c.Status(resp.ErrorCode).JSON(resp.ErrorMessage)
		}

		boundary := string(c.Request().Header.MultipartFormBoundary())
		request := execute.ExecuteCommandRequest[*multipart.Form]{
			AI:      existCommandInfo.AI,
			Command: existCommandInfo.Command,
			Payload: formPayload,
		}

		resp, exeErr := execute.ExecuteFormCommand(request, boundary, h.AiDB, h.Logger)

		if exeErr != nil {
			return c.Status(exeErr.ErrorCode).JSON(exeErr)
		}

		for key, value := range resp.Headers {
			c.Response().Header.Add(key, value)
		}

		return c.Status(resp.Status).Send(resp.Raw.Bytes())
	} else {
		var jsonPayload map[string]interface{}

		if err := c.BodyParser(&jsonPayload); err != nil {
			resp := e.NewErrorResponse(e.HttpInternalError, err.Error())
			return c.Status(resp.ErrorCode).JSON(resp.ErrorMessage)
		}

		request := execute.ExecuteCommandRequest[map[string]interface{}]{
			AI:      existCommandInfo.AI,
			Command: existCommandInfo.Command,
			Payload: jsonPayload,
		}

		resp, exeErr := execute.ExecuteJSONCommand(request, h.AiDB, h.Logger)

		if exeErr != nil {
			return c.Status(exeErr.ErrorCode).JSON(exeErr)
		}

		for key, value := range resp.Headers {
			c.Response().Header.Add(key, value)
		}

		return c.Status(resp.Status).Send(resp.Raw.Bytes())
	}
}
