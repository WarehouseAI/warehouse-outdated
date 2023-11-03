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
	"warehouse/src/services/ai/internal/service/command/get"

	"github.com/sirupsen/logrus"
)

type AIProvider interface {
	GetOneBy(map[string]interface{}) (*pg.AI, *db.DBError)
	Update(*pg.AI, map[string]interface{}) *db.DBError
}

func ExecuteFormDataCommand(formData *multipart.Form, command *get.Response, aiProvider AIProvider, logger *logrus.Logger) (*bytes.Buffer, *httputils.ErrorResponse) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range command.Payload.Payload {
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

	headers := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": fmt.Sprintf("%s %s", string(command.AuthScheme), command.ApiKey),
	}

	response, httpErr := httputils.MakeHTTPRequest(command.Payload.URL, string(command.Payload.RequestScheme), headers, url.Values{}, body)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	buffer, httpErr := httputils.DecodeHTTPResponse(response, command.Payload.OutputType)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	if err := updateUsageCount(command.AI, aiProvider); err != nil {
		return nil, err
	}

	return buffer, nil
}

func ExecuteJSONCommand(jsonData map[string]interface{}, command *get.Response, aiProvider AIProvider, logger *logrus.Logger) (*bytes.Buffer, *httputils.ErrorResponse) {
	json, parseErr := json.Marshal(jsonData)

	if parseErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": parseErr.Error()}).Info("Execute JSON-Command")
		return nil, httputils.NewErrorResponse(httputils.InternalError, parseErr.Error())
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("%s %s", string(command.AuthScheme), command.ApiKey),
	}

	body := bytes.NewBuffer(json)
	response, httpErr := httputils.MakeHTTPRequest(command.Payload.URL, string(command.Payload.RequestScheme), headers, url.Values{}, body)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute JSON-Command")
		return nil, httpErr
	}

	buffer, httpErr := httputils.DecodeHTTPResponse(response, command.Payload.OutputType)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	if err := updateUsageCount(command.AI, aiProvider); err != nil {
		return nil, err
	}

	return buffer, nil
}

func updateUsageCount(ai *pg.AI, aiProvider AIProvider) *httputils.ErrorResponse {
	if err := aiProvider.Update(ai, map[string]interface{}{"used": ai.Used + 1}); err != nil {
		return httputils.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}
