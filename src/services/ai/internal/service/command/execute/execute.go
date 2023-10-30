package execute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"time"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type AIProvider interface {
	GetOneBy(map[string]interface{}) (*pg.AI, *db.DBError)
}

func ExecuteFormDataCommand(formData *multipart.Form, command *pg.Command, aiProvider AIProvider, logger *logrus.Logger) (*bytes.Buffer, *httputils.ErrorResponse) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range command.Payload {
		if value == "file" {
			fileHeader := formData.File[key][0]
			file, err := fileHeader.Open()

			if err != nil {
				logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
				return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
			}

			field, _ := writer.CreateFormFile(key, fileHeader.Filename)
			io.Copy(field, file)

			defer file.Close()
		} else {
			writer.WriteField(key, formData.Value[key][0])
		}
	}

	writer.Close()

	ai, dbErr := aiProvider.GetOneBy(map[string]interface{}{"id": command.AI})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Execute JSON-Command")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	headers := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": fmt.Sprintf("%s %s", string(ai.AuthScheme), ai.ApiKey),
	}

	response, httpErr := httputils.MakeHTTPRequest(command.URL, string(command.RequestScheme), headers, url.Values{}, body)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	buffer, httpErr := httputils.DecodeHTTPResponse(response, command.OutputType)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	return buffer, nil
}

func ExecuteJSONCommand(jsonData map[string]interface{}, command *pg.Command, aiProvider AIProvider, logger *logrus.Logger) (*bytes.Buffer, *httputils.ErrorResponse) {
	json, parseErr := json.Marshal(jsonData)

	if parseErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": parseErr.Error()}).Info("Execute JSON-Command")
		return nil, httputils.NewErrorResponse(httputils.InternalError, parseErr.Error())
	}

	// TODO: Убрать двойное обращение к БД и получать апи ключ ИИшки за один раз (Первый раз обращается в http враппере, второй тут)
	ai, dbErr := aiProvider.GetOneBy(map[string]interface{}{"id": command.AI})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Execute JSON-Command")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("%s %s", string(ai.AuthScheme), ai.ApiKey),
	}

	body := bytes.NewBuffer(json)
	response, httpErr := httputils.MakeHTTPRequest(command.URL, string(command.RequestScheme), headers, url.Values{}, body)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute JSON-Command")
		return nil, httpErr
	}

	buffer, httpErr := httputils.DecodeHTTPResponse(response, command.OutputType)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	return buffer, nil
}
