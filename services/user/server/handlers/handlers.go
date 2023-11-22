package handlers

import (
	"warehouseai/user/adapter/broker/mail"
	"warehouseai/user/adapter/grpc/client/ai"
	"warehouseai/user/adapter/grpc/client/auth"
	"warehouseai/user/dataservice/favoritesdata"
	"warehouseai/user/dataservice/userdata"
	e "warehouseai/user/errors"
	m "warehouseai/user/model"
	"warehouseai/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	UserDB       *userdata.Database
	FavoritesDB  *favoritesdata.Database
	Logger       *logrus.Logger
	MailProducer *mail.MailProducer
	AiClient     *ai.AiGrpcClient
	AuthClient   *auth.AuthGrpcClient
}

func (h *Handler) UpdatePersonalDataHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var newPersonalData service.UpdatePersonalDataRequest

	if err := c.BodyParser(&newPersonalData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	updatedUser, err := service.UpdateUserPersonalData(newPersonalData, userId, h.UserDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (h *Handler) UpdateEmailHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var newEmail service.UpdateUserEmailRequest

	if err := c.BodyParser(&newEmail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	if err := service.UpdateUserEmail(newEmail, userId, h.UserDB, h.MailProducer, h.Logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) UpdatePasswordHandler(c *fiber.Ctx) error {
	var request service.UpdateUserPasswordRequest
	user := c.Locals("user").(*m.User)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	if err := service.UpdateUserPassword(request, user, h.UserDB, h.Logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetFavoritesHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	favorites, err := service.GetFavorites(&service.GetFavoritesRequest{UserId: userId}, h.FavoritesDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(favorites)
}

func (h *Handler) AddFavoriteHandler(c *fiber.Ctx) error {
	var request *service.AddFavoriteRequest
	userId := c.Locals("userId").(string)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	if svcErr := service.AddFavorite(userId, request, h.FavoritesDB, h.AiClient, h.Logger); svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) RemoveFavoriteHandler(c *fiber.Ctx) error {
	var request *service.RemoveFavoriteRequest
	userId := c.Locals("userId").(string)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	if svcErr := service.RemoveFavorite(userId, request, h.FavoritesDB, h.Logger); svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.SendStatus(fiber.StatusOK)
}
