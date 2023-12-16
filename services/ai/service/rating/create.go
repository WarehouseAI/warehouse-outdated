package rating

import (
	"time"
	d "warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	"warehouseai/ai/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type SetAiRatingRequest struct {
	AiId string `json:"ai_id"`
	Rate int16  `json:"rate"`
}

func validateSetRatingRequest(request SetAiRatingRequest, logger *logrus.Logger) *e.ErrorResponse {
	if request.Rate < 1 || request.Rate > 5 {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Invalid rate value"}).Info("Get AI rating")
		return e.NewErrorResponse(e.HttpBadRequest, "Invalid rate value, provide value between 1 and 5")
	}

	return nil
}

func SetAiRating(userId string, request SetAiRatingRequest, ratingRepository d.RatingInterface, logger *logrus.Logger) *e.ErrorResponse {
	validateSetRatingRequest(request, logger)

	newRate := model.RatingPerUser{
		UserId: uuid.Must(uuid.FromString(userId)),
		AiId:   uuid.Must(uuid.FromString(request.AiId)),
		Rate:   request.Rate,
	}

	if err := ratingRepository.Add(newRate); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Set AI rating")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}
