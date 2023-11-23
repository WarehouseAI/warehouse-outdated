package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func UploadImage(pic *multipart.FileHeader, picture dataservice.PictureInterface, logger *logrus.Logger) (string, *e.ErrorResponse) {
	picPayload, err := pic.Open()

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return "", e.NewErrorResponse(e.HttpInternalError, "Can't read provided avatar.")
	}

	url, fileErr := picture.UploadFile(picPayload, fmt.Sprintf("avatar.%s%s", uuid.Must(uuid.NewV4()).String(), filepath.Ext(pic.Filename)))

	if fileErr != nil {
		fmt.Println(fileErr)
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": fileErr.Error()}).Info("Register user")
		return "", e.NewErrorResponse(e.HttpInternalError, "Can't upload user avatar.")
	}

	defer picPayload.Close()
	return url, nil
}
