package update

import (
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Username  string `json:"username"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type UserUpdater interface {
	Update(id string, updatedFields interface{}) (*pg.User, error)
}

func UpdateUser(request Request, userId string, userUpdater UserUpdater, logger *logrus.Logger) (*pg.User, error) {
	updatedUser, err := userUpdater.Update(userId, request)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user")
		return nil, dto.BadRequestError
	}

	return updatedUser, nil
}
