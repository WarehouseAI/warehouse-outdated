package register

import (
	"testing"
	e "warehouseai/auth/errors"

	"github.com/stretchr/testify/require"
)

// Valid request
func TestValidateCorrect(t *testing.T) {
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
			name: "Long password",
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
			name: "Short password",
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
			name: "Invalid email",
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
