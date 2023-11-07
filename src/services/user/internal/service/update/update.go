package update

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	"warehouse/gen"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/grpcutils"
	"warehouse/src/internal/utils/httputils"
	"warehouse/src/internal/utils/mailutils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type RemoveFavoriteAiRequest struct {
	ID string `json:"id"`
}

type UpdateFavoriteAiRequest struct {
	FavoriteAi []*pg.AI `json:"favorite_ai"`
}

type UpdateEmailRequest struct {
	Email            string `json:"email"`
	VerificationCode string `json:"-"`
	Verified         bool   `json:"-"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UserUpdater interface {
	RawUpdate(map[string]interface{}, interface{}) (*pg.User, *db.DBError)
	DeleteAssociation(*pg.User, interface{}, string) *db.DBError
}

type AiProvider interface {
	GetById(context.Context, *gen.GetAiByIdMsg) (*gen.AI, *httputils.ErrorResponse)
}

func UpdateUser(request interface{}, key string, value string, userUpdater UserUpdater, logger *logrus.Logger) (*pg.User, *httputils.ErrorResponse) {
	updatedUser, dbErr := userUpdater.RawUpdate(map[string]interface{}{key: value}, request)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return updatedUser, nil
}

func UpdateUserPassword(request UpdatePasswordRequest, user *pg.User, userUpdater UserUpdater, logger *logrus.Logger) *httputils.ErrorResponse {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user password")
		return httputils.NewErrorResponse(httputils.BadRequest, err.Error())
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	request.Password = string(hash)

	if _, dbErr := userUpdater.RawUpdate(map[string]interface{}{"id": user.ID.String()}, request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user password")
		return httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func AddUserFavoriteAi(request *gen.GetAiByIdMsg, user *pg.User, userUpdater UserUpdater, aiProvider AiProvider, logger *logrus.Logger) *httputils.ErrorResponse {
	ai, gwErr := aiProvider.GetById(context.Background(), request)

	if gwErr != nil {
		return gwErr
	}

	newAi := grpcutils.ProtoToAi(ai)
	newFavorites := append(user.FavoriteAi, &newAi)
	updatedFields := UpdateFavoriteAiRequest{newFavorites}
	if _, dbErr := userUpdater.RawUpdate(map[string]interface{}{"id": user.ID.String()}, updatedFields); dbErr != nil {
		return httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func RemoveUserFavoriteAi(request *gen.GetAiByIdMsg, user *pg.User, userUpdater UserUpdater, aiProvider AiProvider, logger *logrus.Logger) *httputils.ErrorResponse {
	aiProto, gwErr := aiProvider.GetById(context.Background(), request)

	if gwErr != nil {
		return gwErr
	}

	ai := grpcutils.ProtoToAi(aiProto)

	if dbErr := userUpdater.DeleteAssociation(user, &ai, "FavoriteAi"); dbErr != nil {
		return httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func UpdateUserEmail(wg *sync.WaitGroup, respch chan *httputils.ErrorResponse, request UpdateEmailRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) {
	if _, dbErr := userUpdater.RawUpdate(map[string]interface{}{"id": userId}, request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update email")
		respch <- httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	} else {
		respch <- nil
	}

	wg.Done()
}

func SendUpdateNotification(wg *sync.WaitGroup, respch chan *httputils.ErrorResponse, logger *logrus.Logger, request UpdateEmailRequest) {
	message := mailutils.NewMessage(mailutils.EmailVerify, request.Email, fmt.Sprintf("%s/api/user/verify/%s", os.Getenv("DOMAIN_HOST"), request.VerificationCode))

	if mailErr := mailutils.SendEmail(message); mailErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": mailErr.Error()}).Info("Update email")
		respch <- httputils.NewErrorResponse(httputils.InternalError, mailErr.Error())
	} else {
		respch <- nil
	}

	wg.Done()
}
