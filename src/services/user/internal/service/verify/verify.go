package verify

import (
	"time"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Verified         bool    `json:"verified"`
	VerificationCode *string `json:"verification_code"`
}

type UserUpdater interface {
	RawUpdate(map[string]interface{}, interface{}) (*pg.User, *db.DBError)
}

func VerifyUserEmail(request Request, user *pg.User, userUpdater UserUpdater, logger *logrus.Logger) *httputils.ErrorResponse {
	if user.Verified {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Already verified"}).Info("Verify user")
		return httputils.NewErrorResponse(httputils.BadRequest, "User already verified")
	}

	if request.VerificationCode != user.VerificationCode {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Invalid verification code"}).Info("Verify user")
		return httputils.NewErrorResponse(httputils.BadRequest, "Invalid verification code")
	}

	request.VerificationCode = nil

	if _, dbErr := userUpdater.RawUpdate(map[string]interface{}{"id": user.ID.String()}, request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Verify user")
		return httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}
