package register

import (
	"context"
	"testing"
	"time"
	aMock "warehouseai/auth/adapter/mocks"
	dMock "warehouseai/auth/dataservice/mocks"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestVerifyValidate(t *testing.T) {
	request := RegisterVerifyRequest{
		Token:  "someToken",
		UserId: "some-uuid",
	}

	err := validateVerifyRequest(request)
	require.Nil(t, err)
}

func TestVerifyValidateError(t *testing.T) {
	request := RegisterVerifyRequest{
		Token:  "",
		UserId: "",
	}

	expErr := &e.ErrorResponse{
		ErrorCode:    e.HttpBadRequest,
		ErrorMessage: "One of the parameters is empty.",
	}

	err := validateVerifyRequest(request)
	require.NotNil(t, err)
	require.Equal(t, expErr, err)
}

func TestRegisterVerify(t *testing.T) {
	ctl := gomock.NewController(t)
	repositoryMock := dMock.NewMockVerificationTokenInterface(ctl)
	grpcMock := aMock.NewMockUserGrpcInterface(ctl)
	log := logrus.New()

	plainTokenPayload := "some-token"
	hashTokenPayload, _ := bcrypt.GenerateFromPassword([]byte(plainTokenPayload), 12)

	existToken := &m.VerificationToken{
		ID:        uuid.Must(uuid.NewV4()),
		UserId:    uuid.Must(uuid.NewV4()).String(),
		Token:     string(hashTokenPayload),
		ExpiresAt: time.Now().Add(time.Minute * 10),
		CreatedAt: time.Now(),
	}

	request := RegisterVerifyRequest{
		Token:  plainTokenPayload,
		UserId: existToken.UserId,
	}

	repositoryMock.EXPECT().Get(map[string]interface{}{"user_id": existToken.UserId}).Return(existToken, nil).Times(1)
	grpcMock.EXPECT().UpdateVerificationStatus(context.Background(), request.UserId).Return(true, nil).Times(1)
	repositoryMock.EXPECT().Delete(map[string]interface{}{"id": existToken.ID}).Return(nil).Times(1)

	resp, err := RegisterVerify(request, grpcMock, repositoryMock, log)

	require.Nil(t, err)
	require.Equal(t, &RegisterVerifyResponse{Verified: true}, resp)
}
