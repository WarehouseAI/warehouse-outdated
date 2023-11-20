package handlers

import (
	"warehouseai/user/adapter/broker/mail"
	"warehouseai/user/adapter/grpc/client/auth"
	"warehouseai/user/dataservice/userdata"
	e "warehouseai/user/errors"
	"warehouseai/user/model"
	"warehouseai/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	DB           *userdata.Database
	Logger       *logrus.Logger
	MailProducer *mail.MailProducer
	AuthClient   *auth.AuthGrpcClient
}

func (h *Handler) UpdatePersonalDataHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var newPersonalData service.UpdatePersonalDataRequest

	if err := c.BodyParser(&newPersonalData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	updatedUser, err := service.UpdateUserPersonalData(newPersonalData, userId, h.DB, h.Logger)

	if err != nil {
		if err.ErrorCode == e.HttpBadRequest {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (h *Handler) UpdateEmailHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var newEmail service.UpdateUserEmailRequest

	if err := c.BodyParser(&newEmail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	if err := service.UpdateUserEmail(newEmail, userId, h.DB, h.MailProducer, h.Logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) UpdatePasswordHandler(c *fiber.Ctx) error {
	var request service.UpdateUserPasswordRequest
	user := c.Locals("user").(*model.User)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(e.NewErrorResponse(e.HttpBadRequest, "Invalid request body."))
	}

	if err := service.UpdateUserPassword(request, user, h.DB, h.Logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}
