package rating

import (
	"warehouseai/ai/dataservice/ratingdata"
	e "warehouseai/ai/errors"
	"warehouseai/ai/service/rating/get"
	"warehouseai/ai/service/rating/set"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	RatingRepository *ratingdata.Database
	Logger           *logrus.Logger
}

func (h *Handler) GetAiRatingHandler(c *fiber.Ctx) error {
	aiId := c.Query("ai_id")

	rating, err := get.GetAIRating(get.GetAIRatingRequest{AiId: aiId}, h.RatingRepository, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(rating)
}

func (h *Handler) SetRatingForAiHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var rate set.SetAiRatingRequest

	if err := c.BodyParser(&rate); err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, "Invalid request body.")
		return c.Status(response.ErrorCode).JSON(response)
	}

	if err := set.SetAiRating(userId, rate, h.RatingRepository, h.Logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err.ErrorMessage)
	}

	return c.SendStatus(fiber.StatusOK)
}
