package set

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

func validateSetRatingRequest(request *SetAiRatingRequest) *e.ErrorResponse {
	if request.Rate < 1 || request.Rate > 5 {
		return e.NewErrorResponse(e.HttpBadRequest, "Invalid rate value, provide value between 1 and 5")
	}

	return nil
}

func SetAiRating(userId string, request SetAiRatingRequest, aiRepository d.AiInterface, ratingRepository d.RatingInterface, logger *logrus.Logger) *e.ErrorResponse {
	if err := validateSetRatingRequest(&request); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Invalid rate value"}).Info("Get AI rating")
		return err
	}

	if _, err := aiRepository.Get(map[string]interface{}{"id": request.AiId}); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Invalid rate value"}).Info("Get AI rating")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	// Если такой рейтинг уже существует, то обновляем существующую оценку
	existRate, err := ratingRepository.Get(map[string]interface{}{"ai_id": request.AiId, "user_id": userId})

	if err != nil && err.ErrorType == e.DbSystem {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Set AI rating")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	if existRate != nil {
		if err := ratingRepository.Update(existRate, request.Rate); err != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Set AI rating")
			return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
		}

		return nil
	}

	newRate := model.AiRate{
		ByUserId: uuid.Must(uuid.FromString(userId)),
		AiId:     uuid.Must(uuid.FromString(request.AiId)),
		Rate:     request.Rate,
	}

	if err := ratingRepository.Add(&newRate); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Set AI rating")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}
