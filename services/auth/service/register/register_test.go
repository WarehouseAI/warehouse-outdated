package register

import (
	"context"
	"testing"
	"time"
	"warehouseai/auth/adapter/grpc/gen"
	aMock "warehouseai/auth/adapter/mocks"
	dMock "warehouseai/auth/dataservice/mocks"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Valid request
func TestRegisterValidate(t *testing.T) {
	request := &RegisterRequest{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Email:     "validmail@mail.com",
		Username:  "Username",
		Password:  "12345678",
		Image:     "",
		ViaGoogle: false,
	}

	err := validateRegisterRequest(request)
	require.Nil(t, err)
}

func TestValidateError(t *testing.T) {
	cases := []struct {
		name          string
		request       *RegisterRequest
		expectedError *e.ErrorResponse
	}{
		{
			name: "long_password",
			request: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Email:     "validmail@mail.com",
				Username:  "Username",
				Password:  "rqrZBhrHzy9tnNTbL9HzPaAYdtnMqVJ4qEQBkrY77bP5GiaceM5op8642FB3DRMGRA9kSsvaa",
				Image:     "",
				ViaGoogle: false,
			},
			expectedError: &e.ErrorResponse{ErrorCode: e.HttpBadRequest, ErrorMessage: "Password is too long"},
		},
		{
			name: "short_password",
			request: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Email:     "validmail@mail.com",
				Username:  "Username",
				Password:  "1234567",
				Image:     "",
				ViaGoogle: false,
			},
			expectedError: &e.ErrorResponse{ErrorCode: e.HttpBadRequest, ErrorMessage: "Password is too short"},
		},
		{
			name: "bad_email",
			request: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Email:     "validmail",
				Username:  "Username",
				Password:  "12345678",
				Image:     "",
				ViaGoogle: false,
			},
			expectedError: &e.ErrorResponse{ErrorCode: e.HttpBadRequest, ErrorMessage: "The provided string is not email"},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateRegisterRequest(tCase.request)

			require.NotNil(t, err)
			require.Equal(t, tCase.expectedError, err)
		})
	}
}

func TestCreateUser(t *testing.T) {
	ctl := gomock.NewController(t)
	userGatewayMock := aMock.NewMockUserGrpcInterface(ctl)

	user := &gen.CreateUserMsg{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Username:  "username",
		Password:  hashPassword("password"),
		Picture:   "",
		Email:     "validemail@mail.com",
		ViaGoogle: false,
	}

	userGatewayMock.EXPECT().Create(context.Background(), user).Return("id", nil).Times(1)
	resp, err := createUser(user, userGatewayMock)

	require.NotEqual(t, "", resp)
	require.Nil(t, err)
}

func TestCreateUserError(t *testing.T) {
	ctl := gomock.NewController(t)
	userGatewayMock := aMock.NewMockUserGrpcInterface(ctl)

	user := &gen.CreateUserMsg{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Username:  "username",
		Password:  hashPassword("password"),
		Picture:   "",
		Email:     "validemail@mail.com",
		ViaGoogle: false,
	}

	cases := []struct {
		name   string
		expErr *e.ErrorResponse
	}{
		{
			name: "already_exists",
			expErr: &e.ErrorResponse{
				ErrorCode:    e.HttpAlreadyExist,
				ErrorMessage: "User already exists.",
			},
		},
		{
			name: "internal_error",
			expErr: &e.ErrorResponse{
				ErrorCode:    e.HttpInternalError,
				ErrorMessage: "Something went wrong.",
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			userGatewayMock.EXPECT().Create(context.Background(), user).Return("", tCase.expErr).Times(1)
			resp, err := createUser(user, userGatewayMock)

			require.Equal(t, "", resp)
			require.NotNil(t, err)
		})
	}
}

func TestCreateVerificationToken(t *testing.T) {
	ctl := gomock.NewController(t)
	repositoryMock := dMock.NewMockVerificationTokenInterface(ctl)

	validationToken := &m.VerificationToken{
		UserId:    "userID",
		Token:     "newToken",
		ExpiresAt: time.Now().Add(time.Minute * 10),
		CreatedAt: time.Now(),
	}

	repositoryMock.EXPECT().Create(validationToken).Return(nil).Times(1)
	err := createVerificationToken(validationToken, repositoryMock)

	require.Nil(t, err)
}

func TestCreateVerificationTokenError(t *testing.T) {
	ctl := gomock.NewController(t)
	repositoryMock := dMock.NewMockVerificationTokenInterface(ctl)

	validationToken := &m.VerificationToken{
		UserId:    "userID",
		Token:     "newToken",
		ExpiresAt: time.Now().Add(time.Minute * 10),
		CreatedAt: time.Now(),
	}

	cases := []struct {
		name   string
		expErr *e.DBError
	}{
		{
			name: "already_exists",
			expErr: &e.DBError{
				ErrorType: e.DbExist,
				Message:   "Token with this key/keys already exists.",
				Payload:   "some payload",
			},
		},
		{
			name: "internal_error",
			expErr: &e.DBError{
				ErrorType: e.DbSystem,
				Message:   "Something went wrong.",
				Payload:   "some payload",
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			repositoryMock.EXPECT().Create(validationToken).Return(tCase.expErr).Times(1)
			err := createVerificationToken(validationToken, repositoryMock)

			require.NotNil(t, err)
		})
	}
}
