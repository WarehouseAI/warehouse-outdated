package login

import (
	"context"
	"testing"
	"time"
	"warehouseai/auth/adapter/grpc/gen"
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

func TestValidateLogin(t *testing.T) {
	req := &LoginRequest{
		Email:    "validemail@mail.com",
		Password: "12345678",
	}

	err := validateLoginRequest(req)

	require.Nil(t, err)
}

func TestValidateLoginError(t *testing.T) {
	cases := []struct {
		name   string
		req    *LoginRequest
		expErr *e.ErrorResponse
	}{
		{
			name: "invalid_email",
			req: &LoginRequest{
				Email:    "invalidemail",
				Password: "12345678",
			},
			expErr: &e.ErrorResponse{ErrorCode: e.HttpBadRequest, ErrorMessage: "Invalid email address"},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateLoginRequest(tCase.req)

			require.NotNil(t, err)
			require.Equal(t, tCase.expErr, err)
		})
	}
}

func TestLogin(t *testing.T) {
	ctl := gomock.NewController(t)

	grpcMock := aMock.NewMockUserGrpcInterface(ctl)
	dbMock := dMock.NewMockSessionInterface(ctl)
	logger := logrus.New()

	request := &LoginRequest{
		Email:    "validemail@mail.com",
		Password: "12345678",
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), 12)
	expUser := &gen.User{
		Id:        uuid.Must(uuid.NewV4()).String(),
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Username:  "Username",
		Password:  string(hash),
		Email:     request.Email,
		ViaGoogle: false,
		Picture:   "",
		Verified:  true,
		Role:      "Base",
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	sessionPayload := m.SessionPayload{
		UserId:    expUser.Id,
		CreatedAt: time.Now(),
	}
	expSession := &m.Session{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Payload: sessionPayload,
		TTL:     24 * time.Hour,
	}

	grpcMock.EXPECT().GetByEmail(context.Background(), request.Email).Return(expUser, nil).Times(1)
	dbMock.EXPECT().Create(context.Background(), expUser.Id).Return(expSession, nil).Times(1)

	resp, session, err := Login(request, grpcMock, dbMock, logger)

	require.NotNil(t, resp)
	require.NotNil(t, session)
	require.Nil(t, err)
	require.Equal(t, &LoginResponse{UserId: expUser.Id}, resp)
	require.IsType(t, &m.Session{}, session)
}

func TestLoginError(t *testing.T) {
	cases := []struct {
		name   string
		req    *LoginRequest
		expErr *e.ErrorResponse
	}{
		{
			name: "user_not_found",
			req: &LoginRequest{
				Email:    "notexistemail@mail.com",
				Password: "12345678",
			},
			expErr: &e.ErrorResponse{
				ErrorCode:    e.HttpBadRequest,
				ErrorMessage: "User is not exist",
			},
		},
	}

	ctl := gomock.NewController(t)

	grpcMock := aMock.NewMockUserGrpcInterface(ctl)
	dbMock := dMock.NewMockSessionInterface(ctl)
	logger := logrus.New()

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			grpcMock.EXPECT().GetByEmail(context.Background(), tCase.req.Email).Return(nil, &e.ErrorResponse{ErrorCode: e.HttpBadRequest, ErrorMessage: "User is not exist"}).Times(1)

			resp, session, err := Login(tCase.req, grpcMock, dbMock, logger)

			require.Nil(t, resp)
			require.Nil(t, session)
			require.NotNil(t, err)
			require.Equal(t, tCase.expErr, err)
		})
	}
}
