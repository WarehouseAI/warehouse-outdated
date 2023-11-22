package ai

import (
	"fmt"
	"strings"
	"warehouseai/ai/adapter/grpc/client/auth"
	"warehouseai/ai/adapter/grpc/client/user"
	"warehouseai/ai/dataservice/aidata"
	e "warehouseai/ai/errors"
	"warehouseai/ai/service/ai"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	DB         *aidata.Database
	Logger     *logrus.Logger
	UserClient *user.UserGrpcClient
	AuthClient *auth.AuthGrpcClient
}

func (h *Handler) CreateAiWithKeyHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var aiData ai.CreateWithKeyRequest

	if err := c.BodyParser(&aiData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	newAi, err := ai.CreateWithOwnKey(&aiData, userId, h.DB, h.Logger)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(newAi)
}

func (h *Handler) CreateAiWithoutKeyHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var aiData ai.CreateWithoutKeyRequest

	if err := c.BodyParser(&aiData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	newAi, err := ai.CreateWithGeneratedKey(&aiData, userId, h.DB, h.Logger)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(newAi)
}

func (h *Handler) GetAIHandler(c *fiber.Ctx) error {
	aiId := c.Query("id")
	sessionId := c.Cookies("sessionId")

	var existAi *ai.GetAiResponse
	var svcErr *e.ErrorResponse

	if sessionId == "" {
		existAi, svcErr = ai.GetByIdPreload(aiId, h.DB, h.Logger)
	} else {
		userId := c.Locals("userId").(string)
		existAi, svcErr = ai.GetByIdPreloadAuthed(userId, aiId, h.DB, h.UserClient, h.Logger)
	}

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.Status(fiber.StatusOK).JSON(existAi)
}

func (h *Handler) SearchHandler(c *fiber.Ctx) error {
	field := c.Query("field")
	value := "%" + c.Query("value") + "%"

	result, svcErr := ai.GetLike(field, value, h.DB, h.Logger)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) GetAisHandler(c *fiber.Ctx) error {
	plainIds := c.Query("id")
	aiIds := strings.Split(plainIds, ",")

	fmt.Println(aiIds)
	existAis, svcErr := ai.GetManyById(aiIds, h.DB, h.Logger)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.Status(fiber.StatusOK).JSON(existAis)
}
