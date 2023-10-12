package execute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type AIProvider interface {
	GetOneBy(key string, value interface{}) (*pg.AI, error)
}

func ExecuteFormDataCommand(formData *multipart.Form, command *pg.Command, aiProvider AIProvider, logger *logrus.Logger) (*bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range command.Payload {
		if value == "file" {
			fileHeader := formData.File[key][0]
			file, err := fileHeader.Open()

			if err != nil {
				logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
				return nil, dto.InternalError
			}

			field, _ := writer.CreateFormFile(key, fileHeader.Filename)
			io.Copy(field, file)

			defer file.Close()
		} else {
			writer.WriteField(key, formData.Value[key][0])
		}
	}

	writer.Close()

	ai, err := aiProvider.GetOneBy("id", command.AI)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute JSON-Command")
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": fmt.Sprintf("%s %s", string(ai.AuthScheme), ai.ApiKey),
	}

	response, err := httputils.MakeHTTPRequest(command.URL, string(command.RequestScheme), headers, url.Values{}, body)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
		return nil, err
	}

	buffer, err := httputils.DecodeHTTPResponse(response, command.OutputType)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
		return nil, dto.InternalError
	}

	return buffer, nil
}

func ExecuteJSONCommand(jsonData map[string]interface{}, command *pg.Command, aiProvider AIProvider, logger *logrus.Logger) (*bytes.Buffer, error) {
	json, err := json.Marshal(jsonData)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute JSON-Command")
		return nil, dto.InternalError
	}

	// TODO: Убрать двойное обращение к БД и получать апи ключ ИИшки за один раз (Первый раз обращается в http враппере, второй тут)
	ai, err := aiProvider.GetOneBy("id", command.AI)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute JSON-Command")
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("%s %s", string(ai.AuthScheme), ai.ApiKey),
	}

	body := bytes.NewBuffer(json)
	response, err := httputils.MakeHTTPRequest(command.URL, string(command.RequestScheme), headers, url.Values{}, body)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute JSON-Command")
		return nil, err
	}

	buffer, err := httputils.DecodeHTTPResponse(response, command.OutputType)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
		return nil, dto.InternalError
	}

	return buffer, nil
}
