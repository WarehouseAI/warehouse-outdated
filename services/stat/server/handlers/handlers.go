package handlers

import (
	"warehouseai/stat/dataservice/statdata"
	"warehouseai/stat/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	DB     *statdata.Database
	Logger *logrus.Logger
}

func (h *Handler) GetNumOfUsersHandler(c *fiber.Ctx) error {
	num, err := service.GetNumOfUsers(h.DB, h.Logger)
	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(num)
}

func (h *Handler) GetNumOfDevelopersHandler(c *fiber.Ctx) error {
	num, err := service.GetNumOfDevelopers(h.DB, h.Logger)
	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(num)
}

func (h *Handler) GetNumOfAiUsesHandler(c *fiber.Ctx) error {
	aiId := uuid.FromStringOrNil(c.Query("id"))

	uses, err := service.GetNumOfAiUses(aiId, h.DB, h.Logger)
	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(uses)
}
