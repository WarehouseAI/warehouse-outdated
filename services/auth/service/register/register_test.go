package register

import (
	"context"
	"testing"
	"warehouseai/auth/adapter/grpc/gen"
	aMock "warehouseai/auth/adapter/mocks"
	dMock "warehouseai/auth/dataservice/mocks"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
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

func TestRegister(t *testing.T) {
	ctl := gomock.NewController(t)

	grpcMock := aMock.NewMockUserGrpcInterface(ctl)
	dbMock := dMock.NewMockVerificationTokenInterface(ctl)
	brokerMock := aMock.NewMockBrokerInterface(ctl)
	logger := logrus.New()

	request := &RegisterRequest{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Username:  "Username",
		Password:  "12345678",
		Image:     "",
		Email:     "validmail@mail.com",
		ViaGoogle: false,
	}

	userId := uuid.Must(uuid.NewV4()).String()

	grpcMock.EXPECT().Create(context.Background(), gomock.AssignableToTypeOf(&gen.CreateUserMsg{})).Return(userId, nil).Times(1)
	dbMock.EXPECT().Create(gomock.AssignableToTypeOf(&m.VerificationToken{})).Return(nil).Times(1)
	brokerMock.EXPECT().SendEmail(gomock.AssignableToTypeOf(m.Email{})).Return(nil).Times(1)

	resp, err := Register(request, grpcMock, dbMock, brokerMock, logger)

	require.Nil(t, err)
	require.Equal(t, &RegisterResponse{UserId: userId}, resp)
}

func TestRegisterGrpcError(t *testing.T) {
	cases := []struct {
		name   string
		req    *RegisterRequest
		expErr *e.ErrorResponse
	}{
		{
			name: "already_exist",
			req: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Username:  "Username",
				Password:  "12345678",
				Image:     "",
				Email:     "validmail@mail.com",
				ViaGoogle: false,
			},
			expErr: &e.ErrorResponse{
				ErrorCode:    e.HttpAlreadyExist,
				ErrorMessage: "User already exists",
			},
		},
		{
			name: "internal_error",
			req: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Username:  "Username",
				Password:  "12345678",
				Image:     "",
				Email:     "validmail@mail.com",
				ViaGoogle: false,
			},
			expErr: &e.ErrorResponse{
				ErrorCode:    e.HttpInternalError,
				ErrorMessage: "Something went wrong",
			},
		},
	}

	ctl := gomock.NewController(t)

	grpcMock := aMock.NewMockUserGrpcInterface(ctl)
	dbMock := dMock.NewMockVerificationTokenInterface(ctl)
	brokerMock := aMock.NewMockBrokerInterface(ctl)
	logger := logrus.New()

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			grpcMock.EXPECT().Create(context.Background(), gomock.AssignableToTypeOf(&gen.CreateUserMsg{})).Return("", tCase.expErr).Times(1)

			resp, err := Register(tCase.req, grpcMock, dbMock, brokerMock, logger)

			require.NotNil(t, err)
			require.Equal(t, tCase.expErr, err)
			require.Nil(t, resp)
		})
	}
}

func TestRegisterDbError(t *testing.T) {
	cases := []struct {
		name   string
		req    *RegisterRequest
		expErr *e.DBError
	}{
		{
			name: "already_exist",
			req: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Username:  "Username",
				Password:  "12345678",
				Image:     "",
				Email:     "validmail@mail.com",
				ViaGoogle: false,
			},
			expErr: &e.DBError{
				ErrorType: e.DbExist,
				Message:   "Token with this key/keys already exists.",
				Payload:   "token already exists payload",
			},
		},
		{
			name: "internal_error",
			req: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Username:  "Username",
				Password:  "12345678",
				Image:     "",
				Email:     "validmail@mail.com",
				ViaGoogle: false,
			},
			expErr: &e.DBError{
				ErrorType: e.DbSystem,
				Message:   "Something went wrong.",
				Payload:   "internal error payload",
			},
		},
	}

	ctl := gomock.NewController(t)

	grpcMock := aMock.NewMockUserGrpcInterface(ctl)
	dbMock := dMock.NewMockVerificationTokenInterface(ctl)
	brokerMock := aMock.NewMockBrokerInterface(ctl)
	logger := logrus.New()

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			grpcMock.EXPECT().Create(context.Background(), gomock.AssignableToTypeOf(&gen.CreateUserMsg{})).Return("id", nil).Times(1)
			dbMock.EXPECT().Create(gomock.AssignableToTypeOf(&m.VerificationToken{})).Return(tCase.expErr).Times(1)
			brokerMock.EXPECT().SendTokenReject("id").Return(nil).Times(1)

			resp, err := Register(tCase.req, grpcMock, dbMock, brokerMock, logger)

			require.NotNil(t, err)
			require.Equal(t, e.NewErrorResponseFromDBError(tCase.expErr.ErrorType, tCase.expErr.Message), err)
			require.Nil(t, resp)
		})
	}
}
