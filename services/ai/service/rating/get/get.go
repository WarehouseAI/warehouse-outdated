package get

import (
	"math"
	"time"
	d "warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"

	"github.com/sirupsen/logrus"
)

type GetAIRatingRequest struct {
	AiId string `json:"ai_id"`
}

type GetAIRatingResponse struct {
	AverageRating float64 `json:"avg_rating"`
	RatingCount   int64   `json:"count_rating"`
}

func GetAIRating(
	request GetAIRatingRequest,
	ratingRepository d.RatingInterface,
	logger *logrus.Logger,
) (*GetAIRatingResponse, *e.ErrorResponse) {
	rating, err := ratingRepository.GetAverageAiRating(request.AiId)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Get AI rating")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	count, err := ratingRepository.GetCountAiRating(request.AiId)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Get AI rating")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return &GetAIRatingResponse{
		AverageRating: math.Round(*rating*100) / 100,
		RatingCount:   *count,
	}, nil
}
