package set

import (
	"testing"
	dMock "warehouseai/ai/dataservice/mocks"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRatingValidate(t *testing.T) {
	request := &SetAiRatingRequest{AiId: uuid.Must(uuid.NewV4()).String(), Rate: 5}

	err := validateSetRatingRequest(request)

	require.Nil(t, err)
}

func TestRatingValidateError(t *testing.T) {
	cases := []struct {
		name          string
		request       *SetAiRatingRequest
		expectedError *e.ErrorResponse
	}{
		{
			name:          "Rate greater than 5.",
			request:       &SetAiRatingRequest{AiId: uuid.Must(uuid.NewV4()).String(), Rate: 6},
			expectedError: e.NewErrorResponse(e.HttpBadRequest, "Invalid rate value, provide value between 1 and 5"),
		},
		{
			name:          "Rate lower than 1.",
			request:       &SetAiRatingRequest{AiId: uuid.Must(uuid.NewV4()).String(), Rate: 0},
			expectedError: e.NewErrorResponse(e.HttpBadRequest, "Invalid rate value, provide value between 1 and 5"),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateSetRatingRequest(tCase.request)

			require.NotNil(t, err)
			require.Equal(t, tCase.expectedError, err)
		})
	}
}

func TestRatingSet(t *testing.T) {
	ctl := gomock.NewController(t)

	aiMock := dMock.NewMockAiInterface(ctl)
	ratingMock := dMock.NewMockRatingInterface(ctl)
	logger := logrus.New()

	userId := uuid.Must(uuid.NewV4()).String()
	request := SetAiRatingRequest{AiId: uuid.Must(uuid.NewV4()).String(), Rate: 5}

	ratingMock.EXPECT().Get(map[string]interface{}{"ai_id": request.AiId, "user_id": userId}).Return(nil, nil).Times(1)
	ratingMock.EXPECT().Add(&m.RatingPerUser{
		UserId: uuid.Must(uuid.FromString(userId)),
		AiId:   uuid.Must(uuid.FromString(request.AiId)),
		Rate:   request.Rate,
	}).Return(nil).Times(1)

	err := SetAiRating(userId, request, aiMock, ratingMock, logger)

	require.Nil(t, err)
}

func TestRatingUpdate(t *testing.T) {
	ctl := gomock.NewController(t)

	aiMock := dMock.NewMockAiInterface(ctl)
	ratingMock := dMock.NewMockRatingInterface(ctl)
	logger := logrus.New()

	userId := uuid.Must(uuid.NewV4())
	request := SetAiRatingRequest{AiId: uuid.Must(uuid.NewV4()).String(), Rate: 4}

	existRating := &m.RatingPerUser{
		ID:     uuid.Must(uuid.NewV4()),
		AiId:   uuid.Must(uuid.FromString(request.AiId)),
		UserId: uuid.Must(uuid.NewV4()),
		Rate:   5,
	}

	ratingMock.EXPECT().Get(map[string]interface{}{"ai_id": request.AiId, "user_id": userId.String()}).Return(existRating, nil).Times(1)
	ratingMock.EXPECT().Update(existRating, request.Rate).Return(nil).Times(1)

	err := SetAiRating(userId.String(), request, aiMock, ratingMock, logger)

	require.Nil(t, err)
}
