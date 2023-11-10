package handler

import (
	errs "warehouseai/user/errors"
	"warehouseai/user/model"
	"warehouseai/user/service"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) UpdatePersonalData(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var newPersonalData service.UpdatePersonalDataRequest

	if err := c.BodyParser(&newPersonalData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs.NewErrorResponse(errs.HttpBadRequest, "Invalid request body."))
	}

	updatedUser, err := service.UpdateUserPersonalData(newPersonalData, userId, h.db, h.logger)

	if err != nil {
		if err.ErrorCode == errs.HttpBadRequest {
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
		return c.Status(fiber.StatusBadRequest).JSON(errs.NewErrorResponse(errs.HttpBadRequest, "Invalid request body."))
	}

	if err := service.UpdateUserEmail(newEmail, userId, h.db, h.logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) UpdatePasswordHandler(c *fiber.Ctx) error {
	var request service.UpdateUserPasswordRequest
	user := c.Locals("user").(*model.User)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs.NewErrorResponse(errs.HttpBadRequest, "Invalid request body."))
	}

	if err := service.UpdateUserPassword(request, user, h.db, h.logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) UpdateVerificationHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.User)
	verificationCode := c.Params("code")

	if verificationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(errs.NewErrorResponse(errs.HttpBadRequest, "Empty verification code"))
	}

	newVerification := service.UpdateVerificationRequest{Verified: true, VerificationCode: &verificationCode}

	if err := service.UpdateUserVerification(newVerification, user, h.db, h.logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.SendStatus(fiber.StatusOK)
}
