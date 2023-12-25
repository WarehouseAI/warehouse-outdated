package ai

import (
	"fmt"
	"strings"
	"warehouseai/ai/adapter/grpc/client/auth"
	"warehouseai/ai/adapter/grpc/client/user"
	"warehouseai/ai/dataservice/aidata"
	"warehouseai/ai/dataservice/picturedata"
	e "warehouseai/ai/errors"
	"warehouseai/ai/service/ai"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	DB             *aidata.Database
	Logger         *logrus.Logger
	PictureStorage *picturedata.Storage
	UserClient     *user.UserGrpcClient
	AuthClient     *auth.AuthGrpcClient
}

func (h *Handler) CreateAiWithKeyHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	imageUrl := c.Locals("imageUrl").(string)
	form, err := c.MultipartForm()

	if err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	request := ai.CreateWithKeyRequest{
		Name:              form.Value["name"][0],
		AuthHeaderName:    form.Value["auth_header_name"][0],
		AuthHeaderContent: form.Value["auth_header_content"][0],
		Description:       form.Value["description"][0],
		Image:             imageUrl,
	}

	newAi, svcErr := ai.CreateWithOwnKey(&request, userId, h.DB, h.Logger)

	if svcErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(newAi)
}

func (h *Handler) CreateAiWithoutKeyHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	imageUrl := c.Locals("imageUrl").(string)
	form, err := c.MultipartForm()

	if err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	request := ai.CreateWithoutKeyRequest{
		Name:           form.Value["name"][0],
		AuthHeaderName: form.Value["auth_header_name"][0],
		Description:    form.Value["description"][0],
		Image:          imageUrl,
	}

	newAi, svcErr := ai.CreateWithGeneratedKey(&request, userId, h.DB, h.Logger)

	if svcErr != nil {
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
	value := c.Query("value")

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
