package recovery

import (
	"context"
	"time"
	"warehouse/gen"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils"
	"warehouse/src/internal/utils/grpcutils"
	"warehouse/src/internal/utils/httputils"
	"warehouse/src/internal/utils/mailutils"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type ResetAttemptRequest struct {
	Email string `json:"email"`
}

type ResetVerifyResponse struct {
	UserId string `json:"user_id"`
}

type UserProvider interface {
	ResetPassword(context.Context, *gen.ResetPasswordRequest) (*gen.ResetPasswordResponse, *httputils.ErrorResponse)
	GetByEmail(context.Context, *gen.GetUserByEmailMsg) (*gen.User, *httputils.ErrorResponse)
}

type TokenProvider interface {
	Add(*pg.ResetToken) *db.DBError
	DeleteEntity(map[string]interface{}) *db.DBError
	GetOneBy(map[string]interface{}) (*pg.ResetToken, *db.DBError)
}

// TODO: migrate to saga
func PasswordReset(updateInfo *gen.ResetPasswordRequest, resetTokenId string, userProvider UserProvider, tokenProvider TokenProvider, logger *logrus.Logger, ctx context.Context) (*gen.ResetPasswordResponse, *httputils.ErrorResponse) {
	resetToken, dbErr := tokenProvider.GetOneBy(map[string]interface{}{"id": resetTokenId})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Reset password")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	if err := tokenProvider.DeleteEntity(map[string]interface{}{"id": resetToken.ID}); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Reset password")
		return nil, httputils.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(updateInfo.Password), 12)

	if hashErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": hashErr.Error()}).Info("Reset password")
		return nil, httputils.NewErrorResponse(httputils.InternalError, "Error while hashing new password.")
	}

	updateInfo.Password = string(hash)

	resp, gwErr := userProvider.ResetPassword(ctx, updateInfo)

	if gwErr != nil {
		return nil, gwErr
	}

	return resp, nil
}

func VerifyResetCode(verificationCode string, resetTokenId string, tokenProvider TokenProvider, logger *logrus.Logger, ctx context.Context) (*ResetVerifyResponse, *httputils.ErrorResponse) {
	resetToken, err := tokenProvider.GetOneBy(map[string]interface{}{"id": resetTokenId})

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Verify reset token")
		return nil, httputils.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(resetToken.Token), []byte(verificationCode)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Verify reset token")
		return nil, httputils.NewErrorResponse(httputils.BadRequest, "Invalid verification code.")
	}

	return &ResetVerifyResponse{UserId: resetToken.UserId.String()}, nil
}

func SendResetEmail(request ResetAttemptRequest, tokenCreator TokenProvider, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.ResetToken, *httputils.ErrorResponse) {
	protoUser, err := userProvider.GetByEmail(ctx, &gen.GetUserByEmailMsg{Email: request.Email})

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.ErrorMessage}).Info("Send reset token")
		return nil, err
	}

	user := grpcutils.ProtoToUser(protoUser)

	verificationCode := utils.GenerateCode(8)
	hash, _ := bcrypt.GenerateFromPassword([]byte(verificationCode), 12)

	newResetToken := &pg.ResetToken{
		ID:        uuid.Must(uuid.NewV4()),
		UserId:    user.ID,
		Token:     string(hash),
		ExpiresAt: time.Now().Add(time.Minute * 5),
		CreatedAt: time.Now(),
	}

	if err := tokenCreator.Add(newResetToken); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Send reset token")
		return nil, httputils.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	message := mailutils.NewMessage(mailutils.PasswordRecovery, request.Email, verificationCode)

	if err := mailutils.SendEmail(message); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Send reset token")
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	return newResetToken, nil
}
