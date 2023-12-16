package rating

import (
	"time"
	d "warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/sirupsen/logrus"
)

type GetAIRatingRequest struct {
	AiId string `json:"ai_id"`
}

func GetAIRating(request GetAIRatingRequest, ratingRepository d.RatingInterface, logger *logrus.Logger) (*[]m.RatingPerUser, *e.ErrorResponse) {
	rating, err := ratingRepository.GetAiRating(request.AiId)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Get AI rating")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return rating, nil
}
